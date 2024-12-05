package processors

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
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
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
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
	timeouts v1alpha1.DefaultTimeouts,
) StepProcessor {
	if step.Timeouts != nil {
		timeouts = withTimeouts(timeouts, *step.Timeouts)
	}
	return &stepProcessor{
		step:     step,
		report:   report,
		basePath: basePath,
		timeouts: timeouts,
	}
}

type stepProcessor struct {
	step     v1alpha1.TestStep
	report   *model.TestReport
	basePath string
	timeouts v1alpha1.DefaultTimeouts
}

func (p *stepProcessor) Run(ctx context.Context, namespacer namespacer.Namespacer, tc engine.Context) {
	t := testing.FromContext(ctx)
	report := &model.StepReport{
		Name:      p.step.Name,
		StartTime: time.Now(),
	}
	defer func() {
		report.EndTime = time.Now()
		p.report.Add(report)
	}()
	logger := logging.FromContext(ctx)
	if p.step.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*p.step.Compiler))
	}
	contextData := contextData{
		basePath:            p.basePath,
		catch:               p.step.Catch,
		cluster:             p.step.Cluster,
		clusters:            p.step.Clusters,
		deletionPropagation: p.step.DeletionPropagationPolicy,
		skipDelete:          p.step.SkipDelete,
		templating:          p.step.Template,
	}
	tc, err := setupContextAndBindings(ctx, tc, contextData, p.step.Bindings...)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		failer.FailNow(ctx)
	}
	cleaner := cleaner.New(p.timeouts.Cleanup.Duration, tc.DelayBeforeCleanup(), tc.DeletionPropagation())
	t.Cleanup(func() {
		if !cleaner.Empty() || len(p.step.Cleanup) != 0 {
			report := &model.StepReport{
				Name:      fmt.Sprintf("cleanup (%s)", report.Name),
				StartTime: time.Now(),
			}
			defer func() {
				report.EndTime = time.Now()
				p.report.Add(report)
			}()
			logger.Log(logging.Cleanup, logging.BeginStatus, color.BoldFgCyan)
			defer func() {
				logger.Log(logging.Cleanup, logging.EndStatus, color.BoldFgCyan)
			}()
			if !cleaner.Empty() {
				if errs := cleaner.Run(ctx, report); len(errs) != 0 {
					for _, err := range errs {
						logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					}
					failer.Fail(ctx)
				}
			}
			for i, operation := range p.step.Cleanup {
				operationTc := tc
				if operation.Compiler != nil {
					operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
				}
				operations, err := p.finallyOperation(operationTc.Compilers(), i, namespacer, operationTc.Bindings(), operation)
				if err != nil {
					logger.Log(logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.Fail(ctx)
				}
				for _, operation := range operations {
					_, err := operation.execute(ctx, operationTc, report)
					if err != nil {
						failer.Fail(ctx)
					}
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
				operationTc := tc
				if operation.Compiler != nil {
					operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
				}
				operations, err := p.finallyOperation(operationTc.Compilers(), i, namespacer, operationTc.Bindings(), operation)
				if err != nil {
					logger.Log(logging.Finally, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
					failer.Fail(ctx)
				}
				for _, operation := range operations {
					_, err := operation.execute(ctx, operationTc, report)
					if err != nil {
						failer.Fail(ctx)
					}
				}
			}
		}()
	}
	if catch := tc.Catch(); len(catch) != 0 {
		defer func() {
			if t.Failed() {
				logger.Log(logging.Catch, logging.BeginStatus, color.BoldFgCyan)
				defer func() {
					logger.Log(logging.Catch, logging.EndStatus, color.BoldFgCyan)
				}()
				for i, operation := range catch {
					operationTc := tc
					if operation.Compiler != nil {
						operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
					}
					operations, err := p.catchOperation(operationTc.Compilers(), i, namespacer, operationTc.Bindings(), operation)
					if err != nil {
						logger.Log(logging.Catch, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
						failer.Fail(ctx)
					}
					for _, operation := range operations {
						_, err := operation.execute(ctx, operationTc, report)
						if err != nil {
							failer.Fail(ctx)
						}
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
		operationTc := tc
		if operation.Compiler != nil {
			operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
		}
		continueOnError := operation.ContinueOnError != nil && *operation.ContinueOnError
		operations, err := p.tryOperation(operationTc.Compilers(), i, namespacer, operationTc.Bindings(), operation, cleaner)
		if err != nil {
			logger.Log(logging.Try, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			failer.FailNow(ctx)
		}
		for _, operation := range operations {
			outputs, err := operation.execute(ctx, operationTc, report)
			if err != nil {
				if continueOnError {
					failer.Fail(ctx)
				} else {
					failer.FailNow(ctx)
				}
			}
			for k, v := range outputs {
				tc = tc.WithBinding(ctx, k, v)
			}
		}
	}
}

func (p *stepProcessor) tryOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.Operation, cleaner cleaner.CleanerCollector) ([]operation, error) {
	var ops []operation
	if handler.Apply != nil {
		loaded, err := p.applyOperation(compilers, id+1, namespacer, cleaner, bindings, *handler.Apply)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Assert != nil {
		loaded, err := p.assertOperation(compilers, id+1, namespacer, bindings, *handler.Assert)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, p.commandOperation(compilers, id+1, namespacer, *handler.Command))
	} else if handler.Create != nil {
		loaded, err := p.createOperation(compilers, id+1, namespacer, cleaner, bindings, *handler.Create)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Delete != nil {
		loaded, err := p.deleteOperation(compilers, id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Describe != nil {
		ops = append(ops, p.describeOperation(compilers, id+1, namespacer, *handler.Describe))
	} else if handler.Error != nil {
		loaded, err := p.errorOperation(compilers, id+1, namespacer, bindings, *handler.Error)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
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
		ops = append(ops, p.getOperation(compilers, id+1, namespacer, get))
	} else if handler.Get != nil {
		ops = append(ops, p.getOperation(compilers, id+1, namespacer, *handler.Get))
	} else if handler.Patch != nil {
		loaded, err := p.patchOperation(compilers, id+1, namespacer, bindings, *handler.Patch)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.PodLogs != nil {
		ops = append(ops, p.logsOperation(compilers, id+1, namespacer, *handler.PodLogs))
	} else if handler.Proxy != nil {
		ops = append(ops, p.proxyOperation(compilers, id+1, namespacer, *handler.Proxy))
	} else if handler.Script != nil {
		ops = append(ops, p.scriptOperation(compilers, id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, p.sleepOperation(compilers, id+1, *handler.Sleep))
	} else if handler.Update != nil {
		loaded, err := p.updateOperation(compilers, id+1, namespacer, bindings, *handler.Update)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Wait != nil {
		ops = append(ops, p.waitOperation(compilers, id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func (p *stepProcessor) catchOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	if handler.PodLogs != nil {
		ops = append(ops, p.logsOperation(compilers, id+1, namespacer, *handler.PodLogs))
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
		ops = append(ops, p.getOperation(compilers, id+1, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, p.describeOperation(compilers, id+1, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, p.getOperation(compilers, id+1, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := p.deleteOperation(compilers, id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, p.commandOperation(compilers, id+1, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, p.scriptOperation(compilers, id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, p.sleepOperation(compilers, id+1, *handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, p.waitOperation(compilers, id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func (p *stepProcessor) finallyOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	if handler.PodLogs != nil {
		ops = append(ops, p.logsOperation(compilers, id+1, namespacer, *handler.PodLogs))
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
		ops = append(ops, p.getOperation(compilers, id+1, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, p.describeOperation(compilers, id+1, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, p.getOperation(compilers, id+1, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := p.deleteOperation(compilers, id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, p.commandOperation(compilers, id+1, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, p.scriptOperation(compilers, id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, p.sleepOperation(compilers, id+1, *handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, p.waitOperation(compilers, id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func (p *stepProcessor) applyOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Apply) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), compilers, op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypeApply,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				contextData := contextData{
					basePath:   p.basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if err := prepareResource(resource, tc); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opapply.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						getCleanerOrNil(cleaner, tc),
						tc.Templating(),
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

func (p *stepProcessor) assertOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Assert) ([]operation, error) {
	resources, err := p.fileRefOrCheck(context.TODO(), compilers, op.ActionCheckRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypeAssert,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Assert.Duration)
				contextData := contextData{
					basePath:   p.basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					templating: op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opassert.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						tc.Templating(),
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) commandOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.Command) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				op := opcommand.New(
					tc.Compilers(),
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

func (p *stepProcessor) createOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Create) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), compilers, op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypeCreate,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				contextData := contextData{
					basePath:   p.basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if err := prepareResource(resource, tc); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opcreate.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						getCleanerOrNil(cleaner, tc),
						tc.Templating(),
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

func (p *stepProcessor) deleteOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Delete) ([]operation, error) {
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
	resources, err := p.fileRefOrResource(context.TODO(), compilers, ref, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypeDelete,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Delete.Duration)
				contextData := contextData{
					basePath:            p.basePath,
					cluster:             op.Cluster,
					clusters:            op.Clusters,
					deletionPropagation: op.DeletionPropagationPolicy,
					templating:          op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opdelete.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						tc.Templating(),
						tc.DeletionPropagation(),
						op.Expect...,
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) describeOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.Describe) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Describe(ctx, tc.Compilers(), client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					tc.Compilers(),
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

func (p *stepProcessor) errorOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Error) ([]operation, error) {
	resources, err := p.fileRefOrCheck(context.TODO(), compilers, op.ActionCheckRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypeError,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Error.Duration)
				contextData := contextData{
					basePath:   p.basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					templating: op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := operror.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						tc.Templating(),
					)
					return op, timeout, tc, nil
				}
			},
		))
	}
	return ops, nil
}

func (p *stepProcessor) getOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.Get) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Get(ctx, tc.Compilers(), client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					tc.Compilers(),
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

func (p *stepProcessor) logsOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.PodLogs) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Logs(ctx, tc.Compilers(), tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					tc.Compilers(),
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

func (p *stepProcessor) patchOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Patch) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), compilers, op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypePatch,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				contextData := contextData{
					basePath:   p.basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if err := prepareResource(resource, tc); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := oppatch.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						tc.Templating(),
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

func (p *stepProcessor) proxyOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.Proxy) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Proxy(ctx, tc.Compilers(), client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					tc.Compilers(),
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

func (p *stepProcessor) scriptOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.Script) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeScript,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			timeout := timeout.Get(op.Timeout, p.timeouts.Exec.Duration)
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				op := opscript.New(
					tc.Compilers(),
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

func (p *stepProcessor) sleepOperation(_ compilers.Compilers, id int, op v1alpha1.Sleep) operation {
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeSleep,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			return opsleep.New(op), nil, tc, nil
		},
	)
}

func (p *stepProcessor) updateOperation(compilers compilers.Compilers, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Update) ([]operation, error) {
	resources, err := p.fileRefOrResource(context.TODO(), compilers, op.ActionResourceRef, bindings)
	if err != nil {
		return nil, err
	}
	var ops []operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			OperationInfo{
				Id:         id,
				ResourceId: i + 1,
			},
			model.OperationTypeUpdate,
			func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
				timeout := timeout.Get(op.Timeout, p.timeouts.Apply.Duration)
				contextData := contextData{
					basePath:   p.basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
				}
				if tc, err := setupContextAndBindings(ctx, tc, contextData, op.Bindings...); err != nil {
					return nil, nil, tc, err
				} else if err := prepareResource(resource, tc); err != nil {
					return nil, nil, tc, err
				} else if _, client, err := tc.CurrentClusterClient(); err != nil {
					return nil, nil, tc, err
				} else {
					op := opupdate.New(
						tc.Compilers(),
						client,
						resource,
						namespacer,
						tc.Templating(),
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

func (p *stepProcessor) waitOperation(_ compilers.Compilers, id int, namespacer namespacer.Namespacer, op v1alpha1.Wait) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
			// make sure timeout is set to populate the command flag
			op.Timeout = &metav1.Duration{Duration: *timeout.Get(op.Timeout, p.timeouts.Exec.Duration)}
			// shift operation timeout
			timeout := op.Timeout.Duration + 30*time.Second
			contextData := contextData{
				basePath: p.basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
			}
			if tc, err := setupContextAndBindings(ctx, tc, contextData); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				entrypoint, args, err := kubectl.Wait(ctx, tc.Compilers(), client, tc.Bindings(), &op)
				if err != nil {
					return nil, nil, tc, err
				}
				op := opcommand.New(
					tc.Compilers(),
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

func (p *stepProcessor) fileRefOrCheck(ctx context.Context, compilers compilers.Compilers, ref v1alpha1.ActionCheckRef, bindings apis.Bindings) ([]unstructured.Unstructured, error) {
	if ref.Check != nil && ref.Check.Value() != nil {
		if object, ok := ref.Check.Value().(map[string]any); !ok {
			return nil, errors.New("resource must be an object")
		} else {
			return []unstructured.Unstructured{{Object: object}}, nil
		}
	}
	if ref.File != "" {
		ref, err := ref.File.Value(ctx, compilers, bindings)
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

func (p *stepProcessor) fileRefOrResource(ctx context.Context, compilers compilers.Compilers, ref v1alpha1.ActionResourceRef, bindings apis.Bindings) ([]unstructured.Unstructured, error) {
	if ref.Resource != nil {
		return []unstructured.Unstructured{*ref.Resource}, nil
	}
	if ref.File != "" {
		ref, err := ref.File.Value(ctx, compilers, bindings)
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

func prepareResource(resource unstructured.Unstructured, tc engine.Context) error {
	if terminationGrace := tc.TerminationGrace(); terminationGrace != nil {
		seconds := int64(terminationGrace.Seconds())
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

func getCleanerOrNil(cleaner cleaner.CleanerCollector, tc engine.Context) cleaner.CleanerCollector {
	if tc.DryRun() {
		return nil
	}
	if tc.SkipDelete() {
		return nil
	}
	return cleaner
}
