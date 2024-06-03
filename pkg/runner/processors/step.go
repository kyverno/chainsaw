package processors

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/resource"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
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
	opupdate "github.com/kyverno/chainsaw/pkg/runner/operations/update"
	runnertemplate "github.com/kyverno/chainsaw/pkg/runner/template"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

// TODO
// - create if not exists

type StepProcessor interface {
	Run(context.Context, binding.Bindings)
}

func NewStepProcessor(
	config v1alpha1.ConfigurationSpec,
	clusters clusters.Registry,
	namespacer namespacer.Namespacer,
	clock clock.PassiveClock,
	test discovery.Test,
	step v1alpha1.TestStep,
	report *report.StepReport,
) StepProcessor {
	return &stepProcessor{
		config:     config,
		clusters:   clusters,
		namespacer: namespacer,
		clock:      clock,
		test:       test,
		step:       step,
		report:     report,
		timeouts:   config.Timeouts.Combine(test.Spec.Timeouts).Combine(step.Timeouts),
	}
}

type stepProcessor struct {
	config     v1alpha1.ConfigurationSpec
	clusters   clusters.Registry
	namespacer namespacer.Namespacer
	clock      clock.PassiveClock
	test       discovery.Test
	step       v1alpha1.TestStep
	report     *report.StepReport
	timeouts   v1alpha1.Timeouts
}

func (p *stepProcessor) Run(ctx context.Context, bindings binding.Bindings) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
		})
	}
	logger := logging.FromContext(ctx)
	registeredClusters := clusters.Register(p.clusters, p.test.BasePath, p.step.Clusters)
	clusterConfig, clusterClient, err := registeredClusters.Resolve(false, p.test.Spec.Cluster)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	bindings = apibindings.RegisterClusterBindings(ctx, bindings, clusterConfig, clusterClient)
	bindings, err = apibindings.RegisterBindings(ctx, bindings, p.step.Bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	delay := p.config.DelayBeforeCleanup
	if p.test.Spec.DelayBeforeCleanup != nil {
		delay = p.test.Spec.DelayBeforeCleanup
	}
	cleaner := newCleaner(p.namespacer, delay)
	try, err := p.tryOperations(registeredClusters, cleaner)
	if err != nil {
		logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	catch, err := p.catchOperations(registeredClusters)
	if err != nil {
		logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	finally, err := p.finallyOperations(registeredClusters, p.step.Finally...)
	if err != nil {
		logger.Log(logging.Finally, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	cleanup, err := p.finallyOperations(registeredClusters, p.step.Cleanup...)
	if err != nil {
		logger.Log(logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	t.Cleanup(func() {
		if !cleaner.isEmpty() || len(cleanup) != 0 {
			logger.Log(logging.Cleanup, logging.RunStatus, color.BoldFgCyan)
			defer func() {
				logger.Log(logging.Cleanup, logging.DoneStatus, color.BoldFgCyan)
			}()
			cleaner.run(ctx)
			for _, operation := range cleanup {
				operation.execute(ctx, bindings)
			}
		}
	})
	if len(finally) != 0 {
		defer func() {
			logger.Log(logging.Finally, logging.RunStatus, color.BoldFgCyan)
			defer func() {
				logger.Log(logging.Finally, logging.DoneStatus, color.BoldFgCyan)
			}()
			for _, operation := range finally {
				operation.execute(ctx, bindings)
			}
		}()
	}
	if len(catch) != 0 {
		defer func() {
			if t.Failed() {
				logger.Log(logging.Catch, logging.RunStatus, color.BoldFgCyan)
				defer func() {
					logger.Log(logging.Catch, logging.DoneStatus, color.BoldFgCyan)
				}()
				for _, operation := range catch {
					operation.execute(ctx, bindings)
				}
			}
		}()
	}
	logger.Log(logging.Try, logging.RunStatus, color.BoldFgCyan)
	defer func() {
		logger.Log(logging.Try, logging.DoneStatus, color.BoldFgCyan)
	}()
	for _, operation := range try {
		for k, v := range operation.execute(ctx, bindings) {
			bindings = apibindings.RegisterNamedBinding(ctx, bindings, k, v)
		}
	}
}

func (p *stepProcessor) tryOperations(registeredClusters clusters.Registry, cleaner *cleaner) ([]operation, error) {
	var ops []operation
	for i, handler := range p.step.Try {
		register := func(o ...operation) {
			continueOnError := handler.ContinueOnError != nil && *handler.ContinueOnError
			for _, o := range o {
				o.continueOnError = continueOnError
				ops = append(ops, o)
			}
		}
		if handler.Apply != nil {
			loaded, err := p.applyOperation(i+1, registeredClusters, cleaner, *handler.Apply)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Assert != nil {
			loaded, err := p.assertOperation(i+1, registeredClusters, *handler.Assert)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(i+1, registeredClusters, *handler.Command))
		} else if handler.Create != nil {
			loaded, err := p.createOperation(i+1, registeredClusters, cleaner, *handler.Create)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(i+1, registeredClusters, *handler.Delete)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Describe != nil {
			register(p.describeOperation(i+1, registeredClusters, *handler.Describe))
		} else if handler.Error != nil {
			loaded, err := p.errorOperation(i+1, registeredClusters, *handler.Error)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				Cluster:              handler.Events.Cluster,
				Timeout:              handler.Events.Timeout,
				ObjectLabelsSelector: handler.Events.ObjectLabelsSelector,
				Format:               handler.Events.Format,
				ResourceReference: v1alpha1.ResourceReference{
					APIVersion: "v1",
					Kind:       "Event",
				},
			}
			register(p.getOperation(i+1, registeredClusters, get))
		} else if handler.Get != nil {
			register(p.getOperation(i+1, registeredClusters, *handler.Get))
		} else if handler.Patch != nil {
			loaded, err := p.patchOperation(i+1, registeredClusters, *handler.Patch)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.PodLogs != nil {
			register(p.logsOperation(i+1, registeredClusters, *handler.PodLogs))
		} else if handler.Script != nil {
			register(p.scriptOperation(i+1, registeredClusters, *handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(i+1, *handler.Sleep))
		} else if handler.Update != nil {
			loaded, err := p.updateOperation(i+1, registeredClusters, *handler.Update)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Wait != nil {
			register(p.waitOperation(i+1, registeredClusters, *handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) catchOperations(registeredClusters clusters.Registry) ([]operation, error) {
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
	for i, handler := range handlers {
		if handler.PodLogs != nil {
			register(p.logsOperation(i+1, registeredClusters, *handler.PodLogs))
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				Cluster:              handler.Events.Cluster,
				Timeout:              handler.Events.Timeout,
				ObjectLabelsSelector: handler.Events.ObjectLabelsSelector,
				Format:               handler.Events.Format,
				ResourceReference: v1alpha1.ResourceReference{
					APIVersion: "v1",
					Kind:       "Event",
				},
			}
			register(p.getOperation(i+1, registeredClusters, get))
		} else if handler.Describe != nil {
			register(p.describeOperation(i+1, registeredClusters, *handler.Describe))
		} else if handler.Get != nil {
			register(p.getOperation(i+1, registeredClusters, *handler.Get))
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(i+1, registeredClusters, *handler.Delete)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(i+1, registeredClusters, *handler.Command))
		} else if handler.Script != nil {
			register(p.scriptOperation(i+1, registeredClusters, *handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(i+1, *handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(i+1, registeredClusters, *handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) finallyOperations(registeredClusters clusters.Registry, operations ...v1alpha1.Finally) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	for i, handler := range operations {
		if handler.PodLogs != nil {
			register(p.logsOperation(i+1, registeredClusters, *handler.PodLogs))
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				Cluster:              handler.Events.Cluster,
				Timeout:              handler.Events.Timeout,
				ObjectLabelsSelector: handler.Events.ObjectLabelsSelector,
				Format:               handler.Events.Format,
				ResourceReference: v1alpha1.ResourceReference{
					APIVersion: "v1",
					Kind:       "Event",
				},
			}
			register(p.getOperation(i+1, registeredClusters, get))
		} else if handler.Describe != nil {
			register(p.describeOperation(i+1, registeredClusters, *handler.Describe))
		} else if handler.Get != nil {
			register(p.getOperation(i+1, registeredClusters, *handler.Get))
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(i+1, registeredClusters, *handler.Delete)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(i+1, registeredClusters, *handler.Command))
		} else if handler.Script != nil {
			register(p.scriptOperation(i+1, registeredClusters, *handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(i+1, *handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(i+1, registeredClusters, *handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) applyOperation(id int, registeredClusters clusters.Registry, cleaner *cleaner, op v1alpha1.Apply) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Apply "+op.File, report.OperationTypeApply)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	clusterResolver := p.getClusterResolver(clusters.Register(registeredClusters, p.test.BasePath, op.Clusters), op.Cluster)
	for i := range resources {
		resource := resources[i]
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(dryRun)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opapply.New(client, resource, p.namespacer, p.getCleaner(cleaner, dryRun), template, op.Expect, op.Outputs), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) assertOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Assert) ([]operation, error) {
	resources, err := p.fileRefOrCheck(op.FileRefOrCheck)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Assert ", report.OperationTypeAssert)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	clusterResolver := p.getClusterResolver(clusters.Register(registeredClusters, p.test.BasePath, op.Clusters), op.Cluster)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.AssertDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(false)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opassert.New(client, resource, p.namespacer, template), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) commandOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Command) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Command ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			return opcommand.New(op, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
		op.Bindings...,
	)
}

func (p *stepProcessor) createOperation(id int, registeredClusters clusters.Registry, cleaner *cleaner, op v1alpha1.Create) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Create ", report.OperationTypeCreate)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(dryRun)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opcreate.New(client, resource, p.namespacer, p.getCleaner(cleaner, dryRun), template, op.Expect, op.Outputs), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) deleteOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Delete) ([]operation, error) {
	ref := v1alpha1.FileRefOrResource{
		FileRef: v1alpha1.FileRef{
			File: op.File,
		},
	}
	if op.Ref != nil {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(op.Ref.APIVersion)
		resource.SetKind(op.Ref.Kind)
		resource.SetName(op.Ref.Name)
		resource.SetNamespace(op.Ref.Namespace)
		resource.SetLabels(op.Ref.Labels)
		ref.Resource = &resource
	}
	resources, err := p.fileRefOrResource(ref)
	if err != nil {
		return nil, err
	}
	var ops []operation

	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Delete ", report.OperationTypeDelete)
	}
	deletionPropagationPolicy := p.config.DeletionPropagationPolicy
	if op.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *op.DeletionPropagationPolicy
	} else if p.step.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *p.step.DeletionPropagationPolicy
	} else if p.test.Spec.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *p.test.Spec.DeletionPropagationPolicy
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.DeleteDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(false)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opdelete.New(client, resource, p.namespacer, template, deletionPropagationPolicy, op.Expect...), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) describeOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Describe) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Describe ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			cmd, err := kubectl.Describe(client, bindings, &op)
			if err != nil {
				return nil, nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
	)
}

func (p *stepProcessor) errorOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Error) ([]operation, error) {
	resources, err := p.fileRefOrCheck(op.FileRefOrCheck)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Error ", report.OperationTypeCommand)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.ErrorDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(false)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return operror.New(client, resource, p.namespacer, template), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) getOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Get) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Get ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			cmd, err := kubectl.Get(client, bindings, &op)
			if err != nil {
				return nil, nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
	)
}

func (p *stepProcessor) logsOperation(id int, registeredClusters clusters.Registry, op v1alpha1.PodLogs) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Logs ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			cmd, err := kubectl.Logs(bindings, &op)
			if err != nil {
				return nil, nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
	)
}

func (p *stepProcessor) patchOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Patch) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Patch ", report.OperationTypeCreate)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(dryRun)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return oppatch.New(client, resource, p.namespacer, template, op.Expect, op.Outputs), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) scriptOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Script) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Script ", report.OperationTypeScript)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, p.timeouts.ExecDuration()),
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			return opscript.New(op, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
		op.Bindings...,
	)
}

func (p *stepProcessor) sleepOperation(id int, op v1alpha1.Sleep) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Sleep ", report.OperationTypeSleep)
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		nil,
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, nil, nil)
			return opsleep.New(op), bindings, nil
		},
		operationReport,
	)
}

func (p *stepProcessor) updateOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Update) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.FileRefOrResource)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Update ", report.OperationTypeCreate)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Spec.Template, p.config.Template)
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		if err := p.prepareResource(resource); err != nil {
			return nil, err
		}
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, p.timeouts.ApplyDuration()),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(dryRun)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opupdate.New(client, resource, p.namespacer, template, op.Expect, op.Outputs), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) waitOperation(id int, registeredClusters clusters.Registry, op v1alpha1.Wait) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Wait ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	// make sure timeout is set to populate the command flag
	op.Timeout = &metav1.Duration{Duration: *timeout.Get(op.Timeout, p.timeouts.ExecDuration())}
	// shift operation timeout
	timeout := op.Timeout.Duration + 30*time.Second
	registeredClusters = clusters.Register(registeredClusters, p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		&timeout,
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			cmd, err := kubectl.Wait(client, bindings, &op)
			if err != nil {
				return nil, nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
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

func (p *stepProcessor) getClusterResolver(registeredClusters clusters.Registry, opCluster string) func(bool) (*rest.Config, client.Client, error) {
	return func(dryRun bool) (*rest.Config, client.Client, error) {
		return registeredClusters.Resolve(dryRun, opCluster, p.step.Cluster, p.test.Spec.Cluster)
	}
}

func (p *stepProcessor) getCleaner(cleaner *cleaner, dryRun bool) cleanup.Cleaner {
	if dryRun {
		return nil
	}
	if cleanup.Skip(p.config.SkipDelete, p.test.Spec.SkipDelete, p.step.TestStepSpec.SkipDelete) {
		return nil
	}
	return func(obj unstructured.Unstructured, c client.Client) {
		cleaner.addObject(obj, c, timeout.Get(nil, p.timeouts.CleanupDuration()))
	}
}
