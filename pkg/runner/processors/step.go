package processors

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/collect"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	opapply "github.com/kyverno/chainsaw/pkg/runner/operations/apply"
	opassert "github.com/kyverno/chainsaw/pkg/runner/operations/assert"
	opcommand "github.com/kyverno/chainsaw/pkg/runner/operations/command"
	opcreate "github.com/kyverno/chainsaw/pkg/runner/operations/create"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	operror "github.com/kyverno/chainsaw/pkg/runner/operations/error"
	opscript "github.com/kyverno/chainsaw/pkg/runner/operations/script"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/clock"
)

// TODO
// - create if not exists

type StepProcessor interface {
	Run(ctx context.Context)
}

func NewStepProcessor(
	config v1alpha1.ConfigurationSpec,
	client client.Client,
	namespacer namespacer.Namespacer,
	clock clock.PassiveClock,
	test discovery.Test,
	step v1alpha1.TestSpecStep,
	stepReport *report.TestSpecStepReport,
) StepProcessor {
	return &stepProcessor{
		config:     config,
		client:     client,
		namespacer: namespacer,
		clock:      clock,
		test:       test,
		step:       step,
		stepReport: stepReport,
	}
}

type stepProcessor struct {
	config     v1alpha1.ConfigurationSpec
	client     client.Client
	namespacer namespacer.Namespacer
	clock      clock.PassiveClock
	test       discovery.Test
	step       v1alpha1.TestSpecStep
	stepReport *report.TestSpecStepReport
}

func (p *stepProcessor) Run(ctx context.Context) {
	t := testing.FromContext(ctx)
	logger := logging.FromContext(ctx)
	try, err := p.tryOperations(ctx, p.step.TestStepSpec.Try...)
	if err != nil {
		logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	catch, err := p.catchOperations(ctx, p.step.TestStepSpec.Catch...)
	if err != nil {
		logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	finally, err := p.finallyOperations(ctx, p.step.TestStepSpec.Finally...)
	if err != nil {
		logger.Log(logging.Finally, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	if len(catch) != 0 {
		defer func() {
			if t.Failed() {
				t.Cleanup(func() {
					logger.Log(logging.Catch, logging.RunStatus, color.BoldFgCyan)
					defer func() {
						logger.Log(logging.Catch, logging.DoneStatus, color.BoldFgCyan)
					}()
					for _, operation := range catch {
						operation.execute(ctx)
					}
				})
			}
		}()
	}
	if len(finally) != 0 {
		defer func() {
			t.Cleanup(func() {
				logger.Log(logging.Finally, logging.RunStatus, color.BoldFgCyan)
				defer func() {
					logger.Log(logging.Finally, logging.DoneStatus, color.BoldFgCyan)
				}()
				for _, operation := range finally {
					operation.execute(ctx)
				}
			})
		}()
	}
	logger.Log(logging.Try, logging.RunStatus, color.BoldFgCyan)
	defer func() {
		logger.Log(logging.Try, logging.DoneStatus, color.BoldFgCyan)
	}()
	for _, operation := range try {
		operation.execute(ctx)
	}
}

func (p *stepProcessor) tryOperations(ctx context.Context, handlers ...v1alpha1.Operation) ([]operation, error) {
	var ops []operation
	for _, handler := range handlers {
		register := func(o ...operation) {
			continueOnError := handler.ContinueOnError != nil && *handler.ContinueOnError
			for _, o := range o {
				o.continueOnError = continueOnError
				ops = append(ops, o)
			}
		}
		if handler.Apply != nil {
			loaded, err := p.applyOperation(ctx, *handler.Apply, handler.Timeout)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Assert != nil {
			loaded, err := p.assertOperation(ctx, *handler.Assert, handler.Timeout)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(ctx, *handler.Command, handler.Timeout))
		} else if handler.Script != nil {
			register(p.scriptOperation(ctx, *handler.Script, handler.Timeout))
		} else if handler.Create != nil {
			loaded, err := p.createOperation(ctx, *handler.Create, handler.Timeout)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(ctx, *handler.Delete, handler.Timeout)
			if err != nil {
				return nil, err
			}
			register(*loaded)
		} else if handler.Error != nil {
			loaded, err := p.errorOperation(ctx, *handler.Error, handler.Timeout)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) catchOperations(ctx context.Context, handlers ...v1alpha1.Catch) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	for _, handler := range handlers {
		if handler.PodLogs != nil {
			cmd, err := collect.PodLogs(handler.PodLogs)
			if err != nil {
				return nil, err
			}
			register(p.commandOperation(ctx, *cmd, handler.Timeout))
		} else if handler.Events != nil {
			cmd, err := collect.Events(handler.Events)
			if err != nil {
				return nil, err
			}
			register(p.commandOperation(ctx, *cmd, handler.Timeout))
		} else if handler.Command != nil {
			register(p.commandOperation(ctx, *handler.Command, handler.Timeout))
		} else if handler.Script != nil {
			register(p.scriptOperation(ctx, *handler.Script, handler.Timeout))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) finallyOperations(ctx context.Context, handlers ...v1alpha1.Finally) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	for _, handler := range handlers {
		if handler.PodLogs != nil {
			cmd, err := collect.PodLogs(handler.PodLogs)
			if err != nil {
				return nil, err
			}
			register(p.commandOperation(ctx, *cmd, handler.Timeout))
		} else if handler.Events != nil {
			cmd, err := collect.Events(handler.Events)
			if err != nil {
				return nil, err
			}
			register(p.commandOperation(ctx, *cmd, handler.Timeout))
		} else if handler.Command != nil {
			register(p.commandOperation(ctx, *handler.Command, handler.Timeout))
		} else if handler.Script != nil {
			register(p.scriptOperation(ctx, *handler.Script, handler.Timeout))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) applyOperation(ctx context.Context, op v1alpha1.Apply, to *metav1.Duration) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	operationReport := report.NewOperation("Apply "+op.File, report.OperationTypeApply)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	for _, resource := range resources {
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, operation{
			timeout:   timeout.Get(to, timeouts.ApplyDuration()),
			operation: opapply.New(p.getClient(dryRun), resource, p.namespacer, p.getCleaner(ctx, dryRun), op.Expect...),
		})
	}
	return ops, nil
}

func (p *stepProcessor) assertOperation(ctx context.Context, op v1alpha1.Assert, to *metav1.Duration) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	operationReport := report.NewOperation("Assert ", report.OperationTypeAssert)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	for _, resource := range resources {
		ops = append(ops, operation{
			timeout:         timeout.Get(to, timeouts.AssertDuration()),
			operation:       opassert.New(p.client, resource, p.namespacer),
			operationReport: operationReport,
		})
	}
	return ops, nil
}

func (p *stepProcessor) commandOperation(ctx context.Context, exec v1alpha1.Command, to *metav1.Duration) operation {
	operationReport := report.NewOperation("Command ", report.OperationTypeCommand)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	return operation{
		timeout:         timeout.Get(to, timeouts.ExecDuration()),
		operation:       opcommand.New(exec, p.test.BasePath, p.namespacer.GetNamespace()),
		operationReport: operationReport,
	}
}

func (p *stepProcessor) createOperation(ctx context.Context, op v1alpha1.Create, to *metav1.Duration) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	operationReport := report.NewOperation("Create ", report.OperationTypeCreate)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	for _, resource := range resources {
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, operation{
			timeout:   timeout.Get(to, timeouts.ApplyDuration()),
			operation: opcreate.New(p.getClient(dryRun), resource, p.namespacer, p.getCleaner(ctx, dryRun), op.Expect...),
		})
	}
	return ops, nil
}

func (p *stepProcessor) deleteOperation(ctx context.Context, op v1alpha1.Delete, to *metav1.Duration) (*operation, error) {
	var resource unstructured.Unstructured
	resource.SetAPIVersion(op.APIVersion)
	resource.SetKind(op.Kind)
	resource.SetName(op.Name)
	resource.SetNamespace(op.Namespace)
	resource.SetLabels(op.Labels)
	operationReport := report.NewOperation("Delete ", report.OperationTypeDelete)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	return &operation{
		timeout:         timeout.Get(to, timeouts.DeleteDuration()),
		operation:       opdelete.New(p.client, resource, p.namespacer, op.Expect...),
		operationReport: operationReport,
	}, nil
}

func (p *stepProcessor) errorOperation(ctx context.Context, op v1alpha1.Error, to *metav1.Duration) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	operationReport := report.NewOperation("Error ", report.OperationTypeCommand)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	for _, resource := range resources {
		ops = append(ops, operation{
			timeout:         timeout.Get(to, timeouts.ErrorDuration()),
			operation:       operror.New(p.client, resource, p.namespacer),
			operationReport: operationReport,
		})
	}
	return ops, nil
}

func (p *stepProcessor) scriptOperation(ctx context.Context, exec v1alpha1.Script, to *metav1.Duration) operation {
	operationReport := report.NewOperation("Script ", report.OperationTypeScript)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
	return operation{
		timeout:         timeout.Get(to, timeouts.ExecDuration()),
		operation:       opscript.New(exec, p.test.BasePath, p.namespacer.GetNamespace()),
		operationReport: operationReport,
	}
}

func (p *stepProcessor) fileRefOrResource(ref v1alpha1.FileRefOrResource) ([]unstructured.Unstructured, error) {
	if ref.Resource != nil {
		return []unstructured.Unstructured{*ref.Resource}, nil
	}
	if ref.File != "" {
		url, err := url.ParseRequestURI(ref.File)
		if err != nil {
			return resource.Load(filepath.Join(p.test.BasePath, ref.File))
		} else {
			resp, err := http.Get(url.String())
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var resources []unstructured.Unstructured
			if err := yaml.Unmarshal(body, &resources); err != nil {
				return nil, err
			}
			return resources, nil
		}
	}
	return nil, errors.New("file or resource must be set")
}

func (p *stepProcessor) prepareResource(resource unstructured.Unstructured) error {
	if p.config.ForceTerminationGracePeriod != nil {
		seconds := int64(p.config.ForceTerminationGracePeriod.Duration.Seconds())
		switch resource.GetKind() {
		case "Pod":
			if err := unstructured.SetNestedField(resource.UnstructuredContent(), seconds, "spec", "terminationGracePeriodSeconds"); err != nil {
				return err
			}
		case "Deployment", "StatefulSet", "DaemonSet", "Job":
			if err := unstructured.SetNestedField(resource.UnstructuredContent(), seconds, "spec", "template", "spec", "terminationGracePeriodSeconds"); err != nil {
				return err
			}
		case "CronJob":
			if err := unstructured.SetNestedField(resource.UnstructuredContent(), seconds, "spec", "jobTemplate", "spec", "template", "spec", "terminationGracePeriodSeconds"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *stepProcessor) getClient(dryRun bool) client.Client {
	if !dryRun {
		return p.client
	}
	return client.DryRun(p.client)
}

func (p *stepProcessor) getCleaner(ctx context.Context, dryRun bool) cleanup.Cleaner {
	if dryRun {
		return nil
	}
	var cleaner cleanup.Cleaner
	if !cleanup.Skip(p.config.SkipDelete, p.test.Spec.SkipDelete, p.step.TestStepSpec.SkipDelete) {
		cleaner = func(obj unstructured.Unstructured, c client.Client) {
			t := testing.FromContext(ctx)
			t.Cleanup(func() {
				timeouts := timeout.Combine(p.config.Timeouts, p.test.Spec.Timeouts, p.step.TestStepSpec.Timeouts)
				operation := operation{
					continueOnError: true,
					timeout:         timeouts.DeleteDuration(),
					operation:       opdelete.New(c, obj, p.namespacer),
				}
				operation.execute(ctx)
			})
		}
	}
	return cleaner
}
