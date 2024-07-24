package processors

import (
	"context"
	"errors"
	"net/url"
	"path/filepath"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/loaders/resource"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
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

type StepProcessor interface {
	Run(context.Context, model.TestContext)
}

func NewStepProcessor(
	config model.Configuration,
	namespacer namespacer.Namespacer,
	clock clock.PassiveClock,
	test discovery.Test,
	step v1alpha1.TestStep,
	report *report.StepReport,
) StepProcessor {
	return &stepProcessor{
		config:     config,
		namespacer: namespacer,
		clock:      clock,
		test:       test,
		step:       step,
		report:     report,
	}
}

type stepProcessor struct {
	config     model.Configuration
	namespacer namespacer.Namespacer
	clock      clock.PassiveClock
	test       discovery.Test
	step       v1alpha1.TestStep
	report     *report.StepReport
}

func (p *stepProcessor) Run(ctx context.Context, tc model.TestContext) {
	t := testing.FromContext(ctx)
	if p.report != nil {
		p.report.SetStartTime(time.Now())
		t.Cleanup(func() {
			p.report.SetEndTime(time.Now())
		})
	}
	logger := logging.FromContext(ctx)
	if p.step.Timeouts != nil {
		tc = model.WithTimeouts(ctx, tc, *p.step.Timeouts)
	}
	if p.step.TestStepSpec.SkipDelete != nil {
		tc = tc.WithCleanup(ctx, !*p.step.TestStepSpec.SkipDelete)
	}
	tc = model.WithClusters(ctx, tc, p.test.BasePath, p.step.Clusters)
	if _tc, err := model.UseCluster(ctx, tc, p.test.Test.Spec.Cluster); err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	} else {
		tc = _tc
	}
	if _tc, err := model.WithBindings(ctx, tc, p.step.Bindings...); err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	} else {
		tc = _tc
	}
	timeouts := tc.Timeouts()
	delay := p.config.Cleanup.DelayBeforeCleanup
	if p.test.Test.Spec.DelayBeforeCleanup != nil {
		delay = p.test.Test.Spec.DelayBeforeCleanup
	}
	cleaner := cleaner.New(timeouts.Cleanup, delay)
	try, err := p.tryOperations(tc, getCleanerOrNil(cleaner, tc.Cleanup()))
	if err != nil {
		logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	catch, err := p.catchOperations(tc)
	if err != nil {
		logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	finally, err := p.finallyOperations(tc, p.step.Finally...)
	if err != nil {
		logger.Log(logging.Finally, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	cleanup, err := p.finallyOperations(tc, p.step.Cleanup...)
	if err != nil {
		logger.Log(logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	t.Cleanup(func() {
		if !cleaner.Empty() || len(cleanup) != 0 {
			logger.Log(logging.Cleanup, logging.RunStatus, color.BoldFgCyan)
			defer func() {
				logger.Log(logging.Cleanup, logging.DoneStatus, color.BoldFgCyan)
			}()
			for _, err := range cleaner.Run(ctx) {
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx)
			}
			for _, operation := range cleanup {
				operation.execute(ctx, tc.Bindings())
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
				operation.execute(ctx, tc.Bindings())
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
					operation.execute(ctx, tc.Bindings())
				}
			}
		}()
	}
	logger.Log(logging.Try, logging.RunStatus, color.BoldFgCyan)
	defer func() {
		logger.Log(logging.Try, logging.DoneStatus, color.BoldFgCyan)
	}()
	for _, operation := range try {
		for k, v := range operation.execute(ctx, tc.Bindings()) {
			tc = tc.WithBinding(ctx, k, v)
		}
	}
}

func (p *stepProcessor) tryOperations(tc model.TestContext, cleaner cleaner.CleanerCollector) ([]operation, error) {
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
			loaded, err := p.applyOperation(i+1, tc, cleaner, *handler.Apply)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Assert != nil {
			loaded, err := p.assertOperation(i+1, tc, *handler.Assert)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(i+1, tc, *handler.Command))
		} else if handler.Create != nil {
			loaded, err := p.createOperation(i+1, tc, cleaner, *handler.Create)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(i+1, tc, *handler.Delete)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Describe != nil {
			register(p.describeOperation(i+1, tc, *handler.Describe))
		} else if handler.Error != nil {
			loaded, err := p.errorOperation(i+1, tc, *handler.Error)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				ActionClusters: handler.Events.ActionClusters,
				ActionFormat:   handler.Events.ActionFormat,
				ActionTimeout:  handler.Events.ActionTimeout,
				ActionObject: v1alpha1.ActionObject{
					ObjectType: v1alpha1.ObjectType{
						APIVersion: "v1",
						Kind:       "Event",
					},
					ActionObjectSelector: handler.Events.ActionObjectSelector,
				},
			}
			register(p.getOperation(i+1, tc, get))
		} else if handler.Get != nil {
			register(p.getOperation(i+1, tc, *handler.Get))
		} else if handler.Patch != nil {
			loaded, err := p.patchOperation(i+1, tc, *handler.Patch)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.PodLogs != nil {
			register(p.logsOperation(i+1, tc, *handler.PodLogs))
		} else if handler.Proxy != nil {
			register(p.proxyOperation(i+1, tc, *handler.Proxy))
		} else if handler.Script != nil {
			register(p.scriptOperation(i+1, tc, *handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(i+1, tc, *handler.Sleep))
		} else if handler.Update != nil {
			loaded, err := p.updateOperation(i+1, tc, *handler.Update)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Wait != nil {
			register(p.waitOperation(i+1, tc, *handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) catchOperations(tc model.TestContext) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	var handlers []v1alpha1.CatchFinally
	handlers = append(handlers, p.config.Error.Catch...)
	handlers = append(handlers, p.test.Test.Spec.Catch...)
	handlers = append(handlers, p.step.Catch...)
	for i, handler := range handlers {
		if handler.PodLogs != nil {
			register(p.logsOperation(i+1, tc, *handler.PodLogs))
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				ActionClusters: handler.Events.ActionClusters,
				ActionFormat:   handler.Events.ActionFormat,
				ActionTimeout:  handler.Events.ActionTimeout,
				ActionObject: v1alpha1.ActionObject{
					ObjectType: v1alpha1.ObjectType{
						APIVersion: "v1",
						Kind:       "Event",
					},
					ActionObjectSelector: handler.Events.ActionObjectSelector,
				},
			}
			register(p.getOperation(i+1, tc, get))
		} else if handler.Describe != nil {
			register(p.describeOperation(i+1, tc, *handler.Describe))
		} else if handler.Get != nil {
			register(p.getOperation(i+1, tc, *handler.Get))
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(i+1, tc, *handler.Delete)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(i+1, tc, *handler.Command))
		} else if handler.Script != nil {
			register(p.scriptOperation(i+1, tc, *handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(i+1, tc, *handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(i+1, tc, *handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) finallyOperations(tc model.TestContext, operations ...v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	for i, handler := range operations {
		if handler.PodLogs != nil {
			register(p.logsOperation(i+1, tc, *handler.PodLogs))
		} else if handler.Events != nil {
			get := v1alpha1.Get{
				ActionClusters: handler.Events.ActionClusters,
				ActionFormat:   handler.Events.ActionFormat,
				ActionTimeout:  handler.Events.ActionTimeout,
				ActionObject: v1alpha1.ActionObject{
					ObjectType: v1alpha1.ObjectType{
						APIVersion: "v1",
						Kind:       "Event",
					},
					ActionObjectSelector: handler.Events.ActionObjectSelector,
				},
			}
			register(p.getOperation(i+1, tc, get))
		} else if handler.Describe != nil {
			register(p.describeOperation(i+1, tc, *handler.Describe))
		} else if handler.Get != nil {
			register(p.getOperation(i+1, tc, *handler.Get))
		} else if handler.Delete != nil {
			loaded, err := p.deleteOperation(i+1, tc, *handler.Delete)
			if err != nil {
				return nil, err
			}
			register(loaded...)
		} else if handler.Command != nil {
			register(p.commandOperation(i+1, tc, *handler.Command))
		} else if handler.Script != nil {
			register(p.scriptOperation(i+1, tc, *handler.Script))
		} else if handler.Sleep != nil {
			register(p.sleepOperation(i+1, tc, *handler.Sleep))
		} else if handler.Wait != nil {
			register(p.waitOperation(i+1, tc, *handler.Wait))
		} else {
			return nil, errors.New("no operation found")
		}
	}
	return ops, nil
}

func (p *stepProcessor) applyOperation(id int, tc model.TestContext, cleaner cleaner.CleanerCollector, op v1alpha1.Apply) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.ActionResourceRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Apply "+op.File, report.OperationTypeApply)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	clusterResolver := p.getClusterResolver(clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters), op.Cluster)
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
			timeout.Get(op.Timeout, tc.Timeouts().Apply),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(dryRun)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opapply.New(client, resource, p.namespacer, getCleanerOrNil(cleaner, !dryRun), template, op.Expect, op.Outputs), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) assertOperation(id int, tc model.TestContext, op v1alpha1.Assert) ([]operation, error) {
	resources, err := p.fileRefOrCheck(op.ActionCheckRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Assert ", report.OperationTypeAssert)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	clusterResolver := p.getClusterResolver(clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters), op.Cluster)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, tc.Timeouts().Assert),
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

func (p *stepProcessor) commandOperation(id int, tc model.TestContext, op v1alpha1.Command) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Command ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, tc.Timeouts().Exec),
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

func (p *stepProcessor) createOperation(id int, tc model.TestContext, cleaner cleaner.CleanerCollector, op v1alpha1.Create) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.ActionResourceRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Create ", report.OperationTypeCreate)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
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
			timeout.Get(op.Timeout, tc.Timeouts().Apply),
			func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
				config, client, err := clusterResolver(dryRun)
				if err != nil {
					return nil, nil, err
				}
				bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
				return opcreate.New(client, resource, p.namespacer, getCleanerOrNil(cleaner, !dryRun), template, op.Expect, op.Outputs), bindings, nil
			},
			operationReport,
			op.Bindings...,
		))
	}
	return ops, nil
}

func (p *stepProcessor) deleteOperation(id int, tc model.TestContext, op v1alpha1.Delete) ([]operation, error) {
	ref := v1alpha1.ActionResourceRef{
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
	deletionPropagationPolicy := p.config.Deletion.Propagation
	if op.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *op.DeletionPropagationPolicy
	} else if p.step.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *p.step.DeletionPropagationPolicy
	} else if p.test.Test.Spec.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *p.test.Test.Spec.DeletionPropagationPolicy
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, tc.Timeouts().Delete),
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

func (p *stepProcessor) describeOperation(id int, tc model.TestContext, op v1alpha1.Describe) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Describe ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, tc.Timeouts().Exec),
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

func (p *stepProcessor) errorOperation(id int, tc model.TestContext, op v1alpha1.Error) ([]operation, error) {
	resources, err := p.fileRefOrCheck(op.ActionCheckRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Error ", report.OperationTypeCommand)
	}
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			timeout.Get(op.Timeout, tc.Timeouts().Error),
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

func (p *stepProcessor) getOperation(id int, tc model.TestContext, op v1alpha1.Get) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Get ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, tc.Timeouts().Exec),
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

func (p *stepProcessor) logsOperation(id int, tc model.TestContext, op v1alpha1.PodLogs) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Logs ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, tc.Timeouts().Exec),
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

func (p *stepProcessor) patchOperation(id int, tc model.TestContext, op v1alpha1.Patch) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.ActionResourceRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Patch ", report.OperationTypeCreate)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
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
			timeout.Get(op.Timeout, tc.Timeouts().Apply),
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

func (p *stepProcessor) proxyOperation(id int, tc model.TestContext, op v1alpha1.Proxy) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Proxy ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, tc.Timeouts().Exec),
		func(ctx context.Context, bindings binding.Bindings) (operations.Operation, binding.Bindings, error) {
			config, client, err := clusterResolver(false)
			if err != nil {
				return nil, nil, err
			}
			bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
			cmd, err := kubectl.Proxy(client, bindings, &op)
			if err != nil {
				return nil, nil, err
			}
			return opcommand.New(*cmd, p.test.BasePath, ns, config), bindings, nil
		},
		operationReport,
	)
}

func (p *stepProcessor) scriptOperation(id int, tc model.TestContext, op v1alpha1.Script) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Script ", report.OperationTypeScript)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
	clusterResolver := p.getClusterResolver(registeredClusters, op.Cluster)
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		timeout.Get(op.Timeout, tc.Timeouts().Exec),
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

func (p *stepProcessor) sleepOperation(id int, _ model.TestContext, op v1alpha1.Sleep) operation {
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

func (p *stepProcessor) updateOperation(id int, tc model.TestContext, op v1alpha1.Update) ([]operation, error) {
	resources, err := p.fileRefOrResource(op.ActionResourceRef)
	if err != nil {
		return nil, err
	}
	var ops []operation
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Update ", report.OperationTypeCreate)
	}
	dryRun := op.DryRun != nil && *op.DryRun
	template := runnertemplate.Get(op.Template, p.step.Template, p.test.Test.Spec.Template, &p.config.Templating.Enabled)
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
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
			timeout.Get(op.Timeout, tc.Timeouts().Apply),
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

func (p *stepProcessor) waitOperation(id int, tc model.TestContext, op v1alpha1.Wait) operation {
	var operationReport *report.OperationReport
	if p.report != nil {
		operationReport = p.report.ForOperation("Wait ", report.OperationTypeCommand)
	}
	ns := ""
	if p.namespacer != nil {
		ns = p.namespacer.GetNamespace()
	}
	// make sure timeout is set to populate the command flag
	op.Timeout = &metav1.Duration{Duration: *timeout.Get(op.Timeout, tc.Timeouts().Exec)}
	// shift operation timeout
	timeout := op.Timeout.Duration + 30*time.Second
	registeredClusters := clusters.Register(tc.Clusters(), p.test.BasePath, op.Clusters)
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

func (p *stepProcessor) fileRefOrCheck(ref v1alpha1.ActionCheckRef) ([]unstructured.Unstructured, error) {
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

func (p *stepProcessor) fileRefOrResource(ref v1alpha1.ActionResourceRef) ([]unstructured.Unstructured, error) {
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
	terminationGracePeriod := p.config.Execution.ForceTerminationGracePeriod
	if p.test.Test.Spec.ForceTerminationGracePeriod != nil {
		terminationGracePeriod = p.test.Test.Spec.ForceTerminationGracePeriod
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
		return registeredClusters.Resolve(dryRun, opCluster, p.step.Cluster, p.test.Test.Spec.Cluster)
	}
}

func getCleanerOrNil(cleaner cleaner.CleanerCollector, enabled bool) cleaner.CleanerCollector {
	if !enabled {
		return nil
	}
	return cleaner
}
