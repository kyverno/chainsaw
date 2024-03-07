package processors

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/resource"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/kubectl"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	opapply "github.com/kyverno/chainsaw/pkg/runner/operations/apply"
	opassert "github.com/kyverno/chainsaw/pkg/runner/operations/assert"
	opcommand "github.com/kyverno/chainsaw/pkg/runner/operations/command"
	opcreate "github.com/kyverno/chainsaw/pkg/runner/operations/create"
	opdelete "github.com/kyverno/chainsaw/pkg/runner/operations/delete"
	operror "github.com/kyverno/chainsaw/pkg/runner/operations/error"
	oppatch "github.com/kyverno/chainsaw/pkg/runner/operations/patch"
	opscript "github.com/kyverno/chainsaw/pkg/runner/operations/script"
	opsleep "github.com/kyverno/chainsaw/pkg/runner/operations/sleep"
	runnertemplate "github.com/kyverno/chainsaw/pkg/runner/template"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

// TODO
// - create if not exists

type StepProcessor interface {
	Run(ctx context.Context)
}

func NewStepProcessor(
	config v1alpha1.ConfigurationSpec,
	clusters clusters,
	namespacer namespacer.Namespacer,
	clock clock.PassiveClock,
	test discovery.Test,
	step v1alpha1.TestSpecStep,
	stepReport *report.TestSpecStepReport,
	cleaner *cleaner,
	bindings binding.Bindings,
) StepProcessor {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return &stepProcessor{
		config:     config,
		clusters:   clusters,
		namespacer: namespacer,
		clock:      clock,
		test:       test,
		step:       step,
		stepReport: stepReport,
		cleaner:    cleaner,
		bindings:   bindings,
		timeouts:   config.Timeouts.Combine(test.Spec.Timeouts).Combine(step.Timeouts),
	}
}

type stepProcessor struct {
	config     v1alpha1.ConfigurationSpec
	clusters   clusters
	namespacer namespacer.Namespacer
	clock      clock.PassiveClock
	test       discovery.Test
	step       v1alpha1.TestSpecStep
	stepReport *report.TestSpecStepReport
	cleaner    *cleaner
	bindings   binding.Bindings
	timeouts   v1alpha1.Timeouts
}

func (p *stepProcessor) Run(ctx context.Context) {
	t := testing.FromContext(ctx)
	logger := logging.FromContext(ctx)
	config, cluster := p.clusters.client(p.step.Cluster, p.test.Spec.Cluster)
	bindings, err := registerBindings(ctx, p.bindings, config, cluster, p.step.Bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	try, err := p.tryOperations()
	if err != nil {
		logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	catch, err := p.catchOperations()
	if err != nil {
		logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		t.FailNow()
	}
	finally, err := p.finallyOperations()
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
						operation.execute(ctx, bindings)
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
					operation.execute(ctx, bindings)
				}
			})
		}()
	}
	logger.Log(logging.Try, logging.RunStatus, color.BoldFgCyan)
	defer func() {
		logger.Log(logging.Try, logging.DoneStatus, color.BoldFgCyan)
	}()
	for _, operation := range try {
		for k, v := range operation.execute(ctx, bindings) {
			bindings = bindings.Register("$"+k, v)
		}
	}
}

func (p *stepProcessor) tryOperations() ([]operation, error) {
	var ops []operation
	for _, handler := range p.step.Try {
		register := func(o ...operation) {
			continueOnError := handler.ContinueOnError != nil && *handler.ContinueOnError
			for _, o := range o {
				o.continueOnError = continueOnError
				ops = append(ops, o)
			}
		}
		if handler.Apply != nil {
			loaded, err := p.applyOperation(*handler.Apply)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Assert != nil {
			loaded, err := p.assertOperation(*handler.Assert)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(*handler.Command))
		} else if handler.Create != nil {
			loaded, err := p.createOperation(*handler.Create)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Delete != nil {
			loaded := p.deleteOperation(*handler.Delete)
			register(loaded)
		} else if handler.Error != nil {
			loaded, err := p.errorOperation(*handler.Error)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Patch != nil {
			loaded, err := p.patchOperation(*handler.Patch)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Script != nil {
			register(p.scriptOperation(*handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(*handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(*handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) catchOperations() ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	var handlers []v1alpha1.Catch
	handlers = append(handlers, p.config.Catch...)
	handlers = append(handlers, p.test.Spec.Catch...)
	handlers = append(handlers, p.step.Catch...)
	for _, handler := range handlers {
		if handler.PodLogs != nil {
			register(p.logsOperation(*handler.PodLogs))
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				Cluster:              handler.Events.Cluster,
				Timeout:              handler.Events.Timeout,
				ObjectLabelsSelector: handler.Events.ObjectLabelsSelector,
				Format:               handler.Events.Format,
				ResourceReference:    v1alpha1.ResourceReference{Resource: "events"},
			}
			register(p.getOperation(get))
		} else if handler.Describe != nil {
			register(p.describeOperation(*handler.Describe))
		} else if handler.Get != nil {
			register(p.getOperation(*handler.Get))
		} else if handler.Delete != nil {
			loaded := p.deleteOperation(*handler.Delete)
			register(loaded)
		} else if handler.Command != nil {
			register(p.commandOperation(*handler.Command))
		} else if handler.Script != nil {
			register(p.scriptOperation(*handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(*handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(*handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) finallyOperations() ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	for _, handler := range p.step.Finally {
		if handler.PodLogs != nil {
			register(p.logsOperation(*handler.PodLogs))
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				Cluster:              handler.Events.Cluster,
				Timeout:              handler.Events.Timeout,
				ObjectLabelsSelector: handler.Events.ObjectLabelsSelector,
				Format:               handler.Events.Format,
				ResourceReference:    v1alpha1.ResourceReference{Resource: "events"},
			}
			register(p.getOperation(get))
		} else if handler.Describe != nil {
			register(p.describeOperation(*handler.Describe))
		} else if handler.Get != nil {
			register(p.getOperation(*handler.Get))
		} else if handler.Delete != nil {
			loaded := p.deleteOperation(*handler.Delete)
			register(loaded)
		} else if handler.Command != nil {
			register(p.commandOperation(*handler.Command))
		} else if handler.Script != nil {
			register(p.scriptOperation(*handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(*handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(*handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) applyOperation(op v1alpha1.Apply) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Apply "+op.File, report.OperationTypeApply)
		p.stepReport.AddOperation(operationReport)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	config, cluster := p.getClient(op.Cluster, dryRun)
	for _, resource := range resources {
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			opapply.New(cluster, resource, p.namespacer, p.getCleaner(dryRun), template, op.Expect, op.Outputs),
			operationReport,
			config,
			cluster,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) assertOperation(op v1alpha1.Assert) ([]operation, error) {
	resources, err := p.fileRefOrCheck(op.FileRefOrCheck)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Assert ", report.OperationTypeAssert)
		p.stepReport.AddOperation(operationReport)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	for _, resource := range resources {
		ops = append(ops, newOperation(
			false,
			timeout.Get(op.Timeout, p.timeouts.AssertDuration()),
			opassert.New(cluster, resource, p.namespacer, template),
			operationReport,
			config,
			cluster,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) commandOperation(op v1alpha1.Command) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Command ", report.OperationTypeCommand)
		p.stepReport.AddOperation(operationReport)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		opcommand.New(op, p.test.BasePath, ns, config),
		operationReport,
		config,
		cluster,
		op.Bindings...,
	)
}

func (p *stepProcessor) createOperation(op v1alpha1.Create) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Create ", report.OperationTypeCreate)
		p.stepReport.AddOperation(operationReport)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	config, cluster := p.getClient(op.Cluster, dryRun)
	for _, resource := range resources {
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			opcreate.New(cluster, resource, p.namespacer, p.getCleaner(dryRun), template, op.Expect, op.Outputs),
			operationReport,
			config,
			cluster,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) deleteOperation(op v1alpha1.Delete) operation {
	var resource unstructured.Unstructured
	resource.SetAPIVersion(op.APIVersion)
	resource.SetKind(op.Kind)
	resource.SetName(op.Name)
	resource.SetNamespace(op.Namespace)
	resource.SetLabels(op.Labels)
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Delete ", report.OperationTypeDelete)
		p.stepReport.AddOperation(operationReport)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.DeleteDuration()),
		opdelete.New(cluster, resource, p.namespacer, template, op.Expect...),
		operationReport,
		config,
		cluster,
		op.Bindings...,
	)
}

func (p *stepProcessor) describeOperation(op v1alpha1.Describe) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Describe ", report.OperationTypeCommand)
		p.stepReport.AddOperation(operationReport)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newLazyOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(_ context.Context, bindings binding.Bindings) (operations.Operation, error) {
			cmd, err := kubectl.Describe(cluster, bindings, &op)
			if err != nil {
				return nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), nil
		},
		operationReport,
		config,
		cluster,
	)
}

func (p *stepProcessor) errorOperation(op v1alpha1.Error) ([]operation, error) {
	resources, err := p.fileRefOrCheck(op.FileRefOrCheck)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Error ", report.OperationTypeCommand)
		p.stepReport.AddOperation(operationReport)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	for _, resource := range resources {
		ops = append(ops, newOperation(
			false,
			timeout.Get(op.Timeout, p.timeouts.ErrorDuration()),
			operror.New(cluster, resource, p.namespacer, template),
			operationReport,
			config,
			cluster,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) getOperation(op v1alpha1.Get) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Get ", report.OperationTypeCommand)
		p.stepReport.AddOperation(operationReport)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newLazyOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(_ context.Context, bindings binding.Bindings) (operations.Operation, error) {
			cmd, err := kubectl.Get(cluster, bindings, &op)
			if err != nil {
				return nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), nil
		},
		operationReport,
		config,
		cluster,
	)
}

func (p *stepProcessor) logsOperation(op v1alpha1.PodLogs) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Logs ", report.OperationTypeCommand)
		p.stepReport.AddOperation(operationReport)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newLazyOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(_ context.Context, bindings binding.Bindings) (operations.Operation, error) {
			cmd, err := kubectl.Logs(bindings, &op)
			if err != nil {
				return nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), nil
		},
		operationReport,
		config,
		cluster,
	)
}

func (p *stepProcessor) patchOperation(op v1alpha1.Patch) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Patch ", report.OperationTypeCreate)
		p.stepReport.AddOperation(operationReport)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	config, cluster := p.getClient(op.Cluster, dryRun)
	for _, resource := range resources {
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			oppatch.New(cluster, resource, p.namespacer, template, op.Expect),
			operationReport,
			config,
			cluster,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) scriptOperation(op v1alpha1.Script) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Script ", report.OperationTypeScript)
		p.stepReport.AddOperation(operationReport)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		opscript.New(op, p.test.BasePath, ns, config),
		operationReport,
		config,
		cluster,
		op.Bindings...,
	)
}

func (p *stepProcessor) sleepOperation(op v1alpha1.Sleep) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Sleep ", report.OperationTypeSleep)
		p.stepReport.AddOperation(operationReport)
	}
	return newOperation(
		false,
		nil,
		opsleep.New(op),
		operationReport,
		nil,
		nil,
	)
}

func (p *stepProcessor) waitOperation(op v1alpha1.Wait) operation {
	var operationReport *report.OperationReport
	if p.stepReport != nil {
		operationReport = report.NewOperation("Wait ", report.OperationTypeCommand)
		p.stepReport.AddOperation(operationReport)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	config, cluster := p.clusters.client(op.Cluster, p.step.Cluster, p.test.Spec.Cluster)
	return newLazyOperation(
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(_ context.Context, bindings binding.Bindings) (operations.Operation, error) {
			cmd, err := kubectl.Wait(cluster, bindings, &op)
			if err != nil {
				return nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), nil
		},
		operationReport,
		config,
		cluster,
	)
}

func (p *stepProcessor) fileRefOrCheck(ref v1alpha1.FileRefOrCheck) ([]unstructured.Unstructured, error) {
	if ref.Check != nil && ref.Check.Value != nil {
		if object, ok := ref.Check.Value.(map[string]any); !ok {
			return nil, errors.New("resource must be an object")
		} else {
			return []unstructured.Unstructured{{Object: object}}, nil
		}
	}
	if ref.File != "" {
		url, err := url.ParseRequestURI(ref.File)
		if err != nil {
			return resource.Load(filepath.Join(p.test.BasePath, ref.File), false)
		} else {
			return resource.LoadFromURI(url, false)
		}
	}
	return nil, errors.New("file or resource must be set")
}

func (p *stepProcessor) fileRefOrResource(ref v1alpha1.FileRefOrResource) ([]unstructured.Unstructured, error) {
	if ref.Resource != nil {
		return []unstructured.Unstructured{*ref.Resource}, nil
	}
	if ref.File != "" {
		url, err := url.ParseRequestURI(ref.File)
		if err != nil {
			return resource.Load(filepath.Join(p.test.BasePath, ref.File), true)
		} else {
			return resource.LoadFromURI(url, true)
		}
	}
	return nil, errors.New("file or resource must be set")
}

func (p *stepProcessor) prepareResource(resource unstructured.Unstructured) error {
	terminationGracePeriod := p.config.ForceTerminationGracePeriod
	if p.test.Spec.ForceTerminationGracePeriod != nil {
		terminationGracePeriod = p.test.Spec.ForceTerminationGracePeriod
	}
	if terminationGracePeriod != nil {
		seconds := int64(terminationGracePeriod.Seconds())
		if seconds != 0 {
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
	}
	return nil
}

func (p *stepProcessor) getClient(opCluster string, dryRun bool) (*rest.Config, client.Client) {
	config, cluster := p.clusters.client(opCluster, p.step.Cluster, p.test.Spec.Cluster)
	if !dryRun {
		return config, cluster
	}
	return config, client.DryRun(cluster)
}

func (p *stepProcessor) getCleaner(dryRun bool) cleanup.Cleaner {
	if dryRun {
		return nil
	}
	if cleanup.Skip(p.config.SkipDelete, p.test.Spec.SkipDelete, p.step.TestStepSpec.SkipDelete) {
		return nil
	}
	return func(obj unstructured.Unstructured, c client.Client) {
		p.cleaner.register(obj, c, timeout.Get(nil, p.timeouts.CleanupDuration()))
	}
}
