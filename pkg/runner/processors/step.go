package processors

import (
	"context"
	"errors"
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
	try, err := p.tryOperations(ctx, p.step.Spec.Try...)
	if err != nil {
		logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	catch, err := p.catchOperations(ctx, p.step.Spec.Catch...)
	if err != nil {
		logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	finally, err := p.finallyOperations(ctx, p.step.Spec.Finally...)
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
	for _, resource := range resources {
		ops = append(ops, operation{
			timeout:   timeout.Get(timeout.DefaultApplyTimeout, p.config.Timeouts.Apply, p.test.Spec.Timeouts.Apply, p.step.Spec.Timeouts.Apply, to),
			operation: opapply.New(p.getClient(dryRun), resource, p.namespacer, p.getCleaner(ctx, dryRun), op.Expect...),
		})
	}
	return ops, nil
}

func (p *stepProcessor) assertOperation(ctx context.Context, op v1alpha1.Assert, to *metav1.Duration) ([]operation, error) {
	resources, err := p.fileRef(op.FileRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	operationReport := report.NewOperation("Assert ", report.OperationTypeAssert)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	for _, resource := range resources {
		ops = append(ops, operation{
			timeout:         timeout.Get(timeout.DefaultAssertTimeout, p.config.Timeouts.Assert, p.test.Spec.Timeouts.Assert, p.step.Spec.Timeouts.Assert, to),
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
	return operation{
		timeout:         timeout.Get(timeout.DefaultExecTimeout, p.config.Timeouts.Exec, p.test.Spec.Timeouts.Exec, p.step.Spec.Timeouts.Exec, to),
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
	for _, resource := range resources {
		ops = append(ops, operation{
			timeout:   timeout.Get(timeout.DefaultApplyTimeout, p.config.Timeouts.Apply, p.test.Spec.Timeouts.Apply, p.step.Spec.Timeouts.Apply, to),
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
	return &operation{
		timeout:         timeout.Get(timeout.DefaultDeleteTimeout, p.config.Timeouts.Delete, p.test.Spec.Timeouts.Delete, p.step.Spec.Timeouts.Delete, to),
		operation:       opdelete.New(p.client, resource, p.namespacer, op.Check),
		operationReport: operationReport,
	}, nil
}

func (p *stepProcessor) errorOperation(ctx context.Context, op v1alpha1.Error, to *metav1.Duration) ([]operation, error) {
	resources, err := p.fileRef(op.FileRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	operationReport := report.NewOperation("Error ", report.OperationTypeCommand)
	if p.stepReport != nil {
		p.stepReport.AddOperation(operationReport)
	}
	for _, resource := range resources {
		ops = append(ops, operation{
			timeout:         timeout.Get(timeout.DefaultErrorTimeout, p.config.Timeouts.Error, p.test.Spec.Timeouts.Error, p.step.Spec.Timeouts.Error, to),
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
	return operation{
		timeout:         timeout.Get(timeout.DefaultExecTimeout, p.config.Timeouts.Exec, p.test.Spec.Timeouts.Exec, p.step.Spec.Timeouts.Exec, to),
		operation:       opscript.New(exec, p.test.BasePath, p.namespacer.GetNamespace()),
		operationReport: operationReport,
	}
}

func (p *stepProcessor) fileRef(ref v1alpha1.FileRef) ([]unstructured.Unstructured, error) {
	if ref.File != "" {
		return resource.Load(filepath.Join(p.test.BasePath, ref.File))
	}
	return nil, errors.New("file must be set")
}

func (p *stepProcessor) fileRefOrResource(ref v1alpha1.FileRefOrResource) ([]unstructured.Unstructured, error) {
	if ref.Resource != nil {
		return []unstructured.Unstructured{*ref.Resource}, nil
	}
	if ref.File != "" {
		return resource.Load(filepath.Join(p.test.BasePath, ref.File))
	}
	return nil, errors.New("file or resource must be set")
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
	if !cleanup.Skip(p.config.SkipDelete, p.test.Spec.SkipDelete, p.step.Spec.SkipDelete) {
		cleaner = func(obj unstructured.Unstructured, c client.Client) {
			t := testing.FromContext(ctx)
			t.Cleanup(func() {
				operation := operation{
					continueOnError: true,
					timeout:         timeout.Get(timeout.DefaultCleanupTimeout, p.config.Timeouts.Cleanup, p.test.Spec.Timeouts.Cleanup, p.step.Spec.Timeouts.Cleanup, nil),
					operation:       opdelete.New(c, obj, p.namespacer, nil),
				}
				operation.execute(ctx)
			})
		}
	}
	return cleaner
}
