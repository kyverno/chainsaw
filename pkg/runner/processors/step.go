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
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/kubectl"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	opapply "github.com/kyverno/chainsaw/pkg/engine/operations/apply"
	opassert "github.com/kyverno/chainsaw/pkg/engine/operations/assert"
	opcommand "github.com/kyverno/chainsaw/pkg/engine/operations/command"
	opcreate "github.com/kyverno/chainsaw/pkg/engine/operations/create"
	opdelete "github.com/kyverno/chainsaw/pkg/engine/operations/delete"
	operror "github.com/kyverno/chainsaw/pkg/engine/operations/error"
	oppatch "github.com/kyverno/chainsaw/pkg/engine/operations/patch"
	opscript "github.com/kyverno/chainsaw/pkg/engine/operations/script"
	opsleep "github.com/kyverno/chainsaw/pkg/engine/operations/sleep"
	opupdate "github.com/kyverno/chainsaw/pkg/engine/operations/update"
	"github.com/kyverno/chainsaw/pkg/loaders/resource"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/timeout"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/pkg/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type StepProcessor interface {
	Run(context.Context, namespacer.Namespacer, engine.Context)
}

func NewStepProcessor(
	step v1alpha1.TestStep,
	report *model.TestReport,
	basePath string,
	delayBeforeCleanup *time.Duration,
	terminationGracePeriod *metav1.Duration,
	timeouts v1alpha1.DefaultTimeouts,
	deletionPropagationPolicy metav1.DeletionPropagation,
	templating bool,
	skipDelete bool,
	catch ...v1alpha1.CatchFinally,
) StepProcessor {
	if step.Timeouts != nil {
		timeouts = withTimeouts(timeouts, *step.Timeouts)
	}
	if step.DeletionPropagationPolicy != nil {
		deletionPropagationPolicy = *step.DeletionPropagationPolicy
	}
	if step.Template != nil {
		templating = *step.Template
	}
	if step.SkipDelete != nil {
		skipDelete = *step.SkipDelete
	}
	catch = append(catch, step.Catch...)
	return &stepProcessor{
		step:                      step,
		report:                    report,
		basePath:                  basePath,
		delayBeforeCleanup:        delayBeforeCleanup,
		terminationGracePeriod:    terminationGracePeriod,
		timeouts:                  timeouts,
		deletionPropagationPolicy: deletionPropagationPolicy,
		templating:                templating,
		skipDelete:                skipDelete,
		catch:                     catch,
	}
}

type stepProcessor struct {
	step                      v1alpha1.TestStep
	report                    *model.TestReport
	basePath                  string
	delayBeforeCleanup        *time.Duration
	terminationGracePeriod    *metav1.Duration
	timeouts                  v1alpha1.DefaultTimeouts
	deletionPropagationPolicy metav1.DeletionPropagation
	templating                bool
	skipDelete                bool
	catch                     []v1alpha1.CatchFinally
}

func (p *stepProcessor) Run(ctx context.Context, namespacer namespacer.Namespacer, tc engine.Context) {
	t := testing.FromContext(ctx)
	report := model.StepReport{
		Name:      p.step.Name,
		StartTime: time.Now(),
	}
	defer func() {
		report.EndTime = time.Now()
		if t.Failed() {
			report.Failed = true
		}
		p.report.Add(report)
	}()
	logger := logging.FromContext(ctx)
	tc, _, err := setupContextData(ctx, tc, contextData{
		basePath: p.basePath,
		bindings: p.step.Bindings,
		cluster:  p.step.Cluster,
		clusters: p.step.Clusters,
	})
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	cleaner := cleaner.New(p.timeouts.Cleanup.Duration, p.delayBeforeCleanup, p.deletionPropagationPolicy)
	t.Cleanup(func() {
		if !cleaner.Empty() || len(p.step.Cleanup) != 0 {
			logger.Log(logging.Cleanup, logging.BeginStatus, color.BoldFgCyan)
			defer func() {
				logger.Log(logging.Cleanup, logging.EndStatus, color.BoldFgCyan)
			}()
			for _, err := range cleaner.Run(ctx) {
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx)
			}
			for i, operation := range p.step.Cleanup {
				operations, err := p.finallyOperation(i, namespacer, tc.Bindings(), operation)
				if err != nil {
					logger.Log(logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.Fail(ctx)
				}
				for _, operation := range operations {
					operation.execute(ctx, tc)
				}
			}
		}
	})
	if len(p.step.Finally) != 0 {
		defer func() {
			logger.Log(logging.Finally, logging.BeginStatus, color.BoldFgCyan)
			defer func() {
				logger.Log(logging.Finally, logging.EndStatus, color.BoldFgCyan)
			}()
			for i, operation := range p.step.Finally {
				operations, err := p.finallyOperation(i, namespacer, tc.Bindings(), operation)
				if err != nil {
					logger.Log(logging.Finally, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.Fail(ctx)
				}
				for _, operation := range operations {
					operation.execute(ctx, tc)
				}
			}
		}()
	}
	if len(p.catch) != 0 {
		defer func() {
			if t.Failed() {
				logger.Log(logging.Catch, logging.BeginStatus, color.BoldFgCyan)
				defer func() {
					logger.Log(logging.Catch, logging.EndStatus, color.BoldFgCyan)
				}()
				for i, operation := range p.catch {
					operations, err := p.catchOperation(i, namespacer, tc.Bindings(), operation)
					if err != nil {
						logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
						failer.Fail(ctx)
					}
					for _, operation := range operations {
						operation.execute(ctx, tc)
					}
				}
			}
		}()
	}
	logger.Log(logging.Try, logging.BeginStatus, color.BoldFgCyan)
	defer func() {
		logger.Log(logging.Try, logging.EndStatus, color.BoldFgCyan)
	}()
	for i, operation := range p.step.Try {
		operations, err := p.tryOperation(i, namespacer, tc.Bindings(), operation, cleaner)
		if err != nil {
			logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		for _, operation := range operations {
			for k, v := range operation.execute(ctx, tc) {
				tc = tc.WithBinding(ctx, k, v)
			}
		}
	}
}

func (p *stepProcessor) tryOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, handler v1alpha1.Operation, cleaner cleaner.CleanerCollector) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		continueOnError := handler.ContinueOnError != nil && *handler.ContinueOnError
		for _, o := range o {
			o.continueOnError = continueOnError
			ops = append(ops, o)
		}
	}
	if handler.Apply != nil {
		loaded, err := p.applyOperation(id+1, namespacer, cleaner, bindings, *handler.Apply)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Assert != nil {
		loaded, err := p.assertOperation(id+1, namespacer, bindings, *handler.Assert)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Command != nil {
		register(p.commandOperation(id+1, namespacer, *handler.Command))
	} else if handler.Create != nil {
		loaded, err := p.createOperation(id+1, namespacer, cleaner, bindings, *handler.Create)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Delete != nil {
		loaded, err := p.deleteOperation(id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Describe != nil {
		register(p.describeOperation(id+1, namespacer, *handler.Describe))
	} else if handler.Error != nil {
		loaded, err := p.errorOperation(id+1, namespacer, bindings, *handler.Error)
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
		register(p.getOperation(id+1, namespacer, get))
	} else if handler.Get != nil {
		register(p.getOperation(id+1, namespacer, *handler.Get))
	} else if handler.Patch != nil {
		loaded, err := p.patchOperation(id+1, namespacer, bindings, *handler.Patch)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.PodLogs != nil {
		register(p.logsOperation(id+1, namespacer, *handler.PodLogs))
	} else if handler.Proxy != nil {
		register(p.proxyOperation(id+1, namespacer, *handler.Proxy))
	} else if handler.Script != nil {
		register(p.scriptOperation(id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		register(p.sleepOperation(id+1, *handler.Sleep))
	} else if handler.Update != nil {
		loaded, err := p.updateOperation(id+1, namespacer, bindings, *handler.Update)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Wait != nil {
		register(p.waitOperation(id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func (p *stepProcessor) catchOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, handler v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	if handler.PodLogs != nil {
		register(p.logsOperation(id+1, namespacer, *handler.PodLogs))
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
		register(p.getOperation(id+1, namespacer, get))
	} else if handler.Describe != nil {
		register(p.describeOperation(id+1, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		register(p.getOperation(id+1, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := p.deleteOperation(id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Command != nil {
		register(p.commandOperation(id+1, namespacer, *handler.Command))
	} else if handler.Script != nil {
		register(p.scriptOperation(id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		register(p.sleepOperation(id+1, *handler.Sleep))
	} else if handler.Wait != nil {
		register(p.waitOperation(id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func (p *stepProcessor) finallyOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, handler v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	register := func(o ...operation) {
		for _, o := range o {
			o.continueOnError = true
			ops = append(ops, o)
		}
	}
	if handler.PodLogs != nil {
		register(p.logsOperation(id+1, namespacer, *handler.PodLogs))
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
		register(p.getOperation(id+1, namespacer, get))
	} else if handler.Describe != nil {
		register(p.describeOperation(id+1, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		register(p.getOperation(id+1, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := p.deleteOperation(id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		register(loaded...)
	} else if handler.Command != nil {
		register(p.commandOperation(id+1, namespacer, *handler.Command))
	} else if handler.Script != nil {
		register(p.scriptOperation(id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		register(p.sleepOperation(id+1, *handler.Sleep))
	} else if handler.Wait != nil {
		register(p.waitOperation(id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func (p *stepProcessor) applyOperation(id int, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings binding.Bindings, op v1alpha1.Apply) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	template := p.getTemplating(op.Template)
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
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
					dryRun:   op.DryRun,
				}); err != nil {
					return nil, nil, tc, err
				} else {
					if _, client, err := tc.CurrentClusterClient(); err != nil {
						return nil, nil, tc, err
					} else {
						op := opapply.New(
							client,
							resource,
							namespacer,
							p.getCleanerOrNil(cleaner, tc),
							template,
							op.Expect,
							op.Outputs,
						)
						return op, timeout, tc, nil
					}
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) assertOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, op v1alpha1.Assert) ([]operation, error) {
	resources, err := p.fileRefOrCheck(context.TODO(), op.ActionCheckRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	template := p.getTemplating(op.Template)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Assert.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
				}); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opassert.New(
						client,
						resource,
						namespacer,
						template,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) commandOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.Command) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: op.Bindings,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				op := opcommand.New(
					op,
					p.basePath,
					ns,
					config,
				)
				return op, timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) createOperation(id int, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings binding.Bindings, op v1alpha1.Create) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	template := p.getTemplating(op.Template)
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
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
					dryRun:   op.DryRun,
				}); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opcreate.New(
						client,
						resource,
						namespacer,
						p.getCleanerOrNil(cleaner, tc),
						template,
						op.Expect,
						op.Outputs,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) deleteOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, op v1alpha1.Delete) ([]operation, error) {
	ref := v1alpha1.ActionResourceRef{
		FileRef: v1alpha1.FileRef{
			File: op.File,
		},
	}
	if op.Ref != nil {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(string(op.Ref.APIVersion))
		resource.SetKind(string(op.Ref.Kind))
		resource.SetName(string(op.Ref.Name))
		resource.SetNamespace(string(op.Ref.Namespace))
		resource.SetLabels(op.Ref.Labels)
		ref.Resource = &resource
	}
	resources, err := p.fileRefOrResource(context.TODO(), ref, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	deletionPropagationPolicy := p.getDeletionPropagationPolicy(op.DeletionPropagationPolicy)
	template := p.getTemplating(op.Template)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Delete.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
				}); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opdelete.New(
						client,
						resource,
						namespacer,
						template,
						deletionPropagationPolicy,
						op.Expect...,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) describeOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.Describe) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: nil,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Describe(ctx, client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					v1alpha1.Command{
						ActionClusters: op.ActionClusters,
						ActionTimeout:  op.ActionTimeout,
						Entrypoint:     entrypoint,
						Args:           args,
					},
					p.basePath,
					ns,
					config,
				)
				return op, timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) errorOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, op v1alpha1.Error) ([]operation, error) {
	resources, err := p.fileRefOrCheck(context.TODO(), op.ActionCheckRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	template := p.getTemplating(op.Template)
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			false,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Error.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
				}); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := operror.New(
						client,
						resource,
						namespacer,
						template,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) getOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.Get) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: nil,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Get(ctx, client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					v1alpha1.Command{
						ActionClusters: op.ActionClusters,
						ActionTimeout:  op.ActionTimeout,
						Entrypoint:     entrypoint,
						Args:           args,
					},
					p.basePath,
					ns,
					config,
				)
				return op, timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) logsOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.PodLogs) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: nil,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Logs(ctx, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					v1alpha1.Command{
						ActionClusters: op.ActionClusters,
						ActionTimeout:  op.ActionTimeout,
						Entrypoint:     entrypoint,
						Args:           args,
					},
					p.basePath,
					ns,
					config,
				)
				return op, timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) patchOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, op v1alpha1.Patch) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	template := p.getTemplating(op.Template)
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
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
					dryRun:   op.DryRun,
				}); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := oppatch.New(
						client,
						resource,
						namespacer,
						template,
						op.Expect,
						op.Outputs,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) proxyOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.Proxy) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: nil,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Proxy(ctx, client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					v1alpha1.Command{
						ActionClusters: op.ActionClusters,
						ActionTimeout:  op.ActionTimeout,
						Entrypoint:     entrypoint,
						Args:           args,
					},
					p.basePath,
					ns,
					config,
				)
				return op, timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) scriptOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.Script) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: op.Bindings,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				op := opscript.New(
					op,
					p.basePath,
					ns,
					config,
				)
				return op, timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) sleepOperation(id int, op v1alpha1.Sleep) operation {
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			return opsleep.New(op), nil, tc, nil
		},
	)
}

func (p *stepProcessor) updateOperation(id int, namespacer namespacer.Namespacer, bindings binding.Bindings, op v1alpha1.Update) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	template := p.getTemplating(op.Template)
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
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				if tc, _, err := setupContextData(ctx, tc, contextData{
					basePath: p.basePath,
					bindings: op.Bindings,
					cluster:  op.Cluster,
					clusters: op.Clusters,
					dryRun:   op.DryRun,
				}); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opupdate.New(
						client,
						resource,
						namespacer,
						template,
						op.Expect,
						op.Outputs,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) waitOperation(id int, namespacer namespacer.Namespacer, op v1alpha1.Wait) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		false,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			// make sure timeout is set to populate the command flag
			op.Timeout = &metav1.Duration{Duration: *timeout.Get(op.Timeout, p.timeouts.Exec.Duration)}
			// shift operation timeout
			timeout := op.Timeout.Duration + 30*time.Second
			if tc, _, err := setupContextData(ctx, tc, contextData{
				basePath: p.basePath,
				bindings: nil,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Wait(ctx, client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					v1alpha1.Command{
						ActionClusters: op.ActionClusters,
						ActionTimeout:  op.ActionTimeout,
						Entrypoint:     entrypoint,
						Args:           args,
					},
					p.basePath,
					ns,
					config,
				)
				return op, &timeout, tc, nil
			}
		},
	)
}

func (p *stepProcessor) fileRefOrCheck(ctx context.Context, ref v1alpha1.ActionCheckRef, bindings binding.Bindings) ([]unstructured.Unstructured, error) {
	if ref.Check != nil && ref.Check.Value != nil {
		if object, ok := ref.Check.Value.(map[string]any); !ok {
			return nil, errors.New("resource must be an object")
		} else {
			return []unstructured.Unstructured{{Object: object}}, nil
		}
	}
	if ref.File != "" {
		ref, err := ref.File.Value(ctx, bindings)
		if err != nil {
			return nil, err
		}
		url, err := url.ParseRequestURI(ref)
		if err != nil {
			return resource.Load(filepath.Join(p.basePath, ref), false)
		} else {
			return resource.LoadFromURI(url, false)
		}
	}
	return nil, errors.New("file or resource must be set")
}

func (p *stepProcessor) fileRefOrResource(ctx context.Context, ref v1alpha1.ActionResourceRef, bindings binding.Bindings) ([]unstructured.Unstructured, error) {
	if ref.Resource != nil {
		return []unstructured.Unstructured{*ref.Resource}, nil
	}
	if ref.File != "" {
		ref, err := ref.File.Value(ctx, bindings)
		if err != nil {
			return nil, err
		}
		url, err := url.ParseRequestURI(ref)
		if err != nil {
			return resource.Load(filepath.Join(p.basePath, ref), true)
		} else {
			return resource.LoadFromURI(url, true)
		}
	}
	return nil, errors.New("file or resource must be set")
}

func (p *stepProcessor) prepareResource(resource unstructured.Unstructured) error {
	if p.terminationGracePeriod != nil {
		seconds := int64(p.terminationGracePeriod.Seconds())
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

func (p *stepProcessor) getCleanerOrNil(cleaner cleaner.CleanerCollector, tc engine.Context) cleaner.CleanerCollector {
	if tc.DryRun() {
		return nil
	}
	if p.skipDelete {
		return nil
	}
	return cleaner
}

func (p *stepProcessor) getTemplating(op *bool) bool {
	if op != nil {
		return *op
	}
	return p.templating
}

func (p *stepProcessor) getDeletionPropagationPolicy(op *metav1.DeletionPropagation) metav1.DeletionPropagation {
	if op != nil {
		return *op
	}
	return p.deletionPropagationPolicy
}
