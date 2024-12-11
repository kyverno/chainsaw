package runner

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
	"github.com/kyverno/chainsaw/pkg/engine/kubectl"
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
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/pkg/ext/output/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
)

func (r *runner) runStep(ctx context.Context, t testing.TTest, basePath string, namespacer namespacer.Namespacer, tc enginecontext.TestContext, step v1alpha1.TestStep, testReport *model.TestReport) bool {
	report := &model.StepReport{
		Name:      step.Name,
		StartTime: time.Now(),
	}
	defer func() {
		report.EndTime = time.Now()
		testReport.Add(report)
	}()
	if step.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*step.Compiler))
	}
	contextData := contextData{
		basePath:            basePath,
		catch:               step.Catch,
		cluster:             step.Cluster,
		clusters:            step.Clusters,
		deletionPropagation: step.DeletionPropagationPolicy,
		skipDelete:          step.SkipDelete,
		templating:          step.Template,
		timeouts:            step.Timeouts,
	}
	tc, err := setupContextAndBindings(tc, contextData, step.Bindings...)
	if err != nil {
		t.Fail()
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		r.onFail()
		return true
	}
	cleaner := cleaner.New(tc.Timeouts().Cleanup.Duration, tc.DelayBeforeCleanup(), tc.DeletionPropagation())
	t.Cleanup(func() {
		if !cleaner.Empty() || len(step.Cleanup) != 0 {
			report := &model.StepReport{
				Name:      fmt.Sprintf("cleanup (%s)", report.Name),
				StartTime: time.Now(),
			}
			defer func() {
				report.EndTime = time.Now()
				testReport.Add(report)
			}()
			logging.Log(ctx, logging.Cleanup, logging.BeginStatus, nil, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.EndStatus, nil, color.BoldFgCyan)
			}()
			if !cleaner.Empty() {
				if errs := cleaner.Run(ctx, report); len(errs) != 0 {
					t.Fail()
					for _, err := range errs {
						logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					}
					r.onFail()
				}
			}
			for i, operation := range step.Cleanup {
				operationTc := tc
				if operation.Compiler != nil {
					operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
				}
				operations, err := finallyOperation(operationTc.Compilers(), basePath, i, namespacer, operationTc.Bindings(), operation)
				if err != nil {
					t.Fail()
					logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
				}
				for _, operation := range operations {
					_, err := operation.execute(ctx, operationTc, report)
					if err != nil {
						t.Fail()
						r.onFail()
					}
				}
			}
		}
	})
	if len(step.Finally) != 0 {
		defer func() {
			logging.Log(ctx, logging.Finally, logging.BeginStatus, nil, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Finally, logging.EndStatus, nil, color.BoldFgCyan)
			}()
			for i, operation := range step.Finally {
				operationTc := tc
				if operation.Compiler != nil {
					operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
				}
				operations, err := finallyOperation(operationTc.Compilers(), basePath, i, namespacer, operationTc.Bindings(), operation)
				if err != nil {
					t.Fail()
					logging.Log(ctx, logging.Finally, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
					r.onFail()
				}
				for _, operation := range operations {
					_, err := operation.execute(ctx, operationTc, report)
					if err != nil {
						t.Fail()
						r.onFail()
					}
				}
			}
		}()
	}
	if catch := tc.Catch(); len(catch) != 0 {
		defer func() {
			if t.Failed() {
				logging.Log(ctx, logging.Catch, logging.BeginStatus, nil, color.BoldFgCyan)
				defer func() {
					logging.Log(ctx, logging.Catch, logging.EndStatus, nil, color.BoldFgCyan)
				}()
				for i, operation := range catch {
					operationTc := tc
					if operation.Compiler != nil {
						operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
					}
					operations, err := catchOperation(operationTc.Compilers(), basePath, i, namespacer, operationTc.Bindings(), operation)
					if err != nil {
						t.Fail()
						logging.Log(ctx, logging.Catch, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
						r.onFail()
					}
					for _, operation := range operations {
						_, err := operation.execute(ctx, operationTc, report)
						if err != nil {
							t.Fail()
							r.onFail()
						}
					}
				}
			}
		}()
	}
	logging.Log(ctx, logging.Try, logging.BeginStatus, nil, color.BoldFgCyan)
	defer func() {
		logging.Log(ctx, logging.Try, logging.EndStatus, nil, color.BoldFgCyan)
	}()
	for i, operation := range step.Try {
		operationTc := tc
		if operation.Compiler != nil {
			operationTc = operationTc.WithDefaultCompiler(string(*operation.Compiler))
		}
		continueOnError := operation.ContinueOnError != nil && *operation.ContinueOnError
		operations, err := tryOperation(operationTc.Compilers(), basePath, i, namespacer, operationTc.Bindings(), operation, cleaner)
		if err != nil {
			t.Fail()
			logging.Log(ctx, logging.Try, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
			r.onFail()
			return true
		}
		for _, operation := range operations {
			outputs, err := operation.execute(ctx, operationTc, report)
			if err != nil {
				t.Fail()
				r.onFail()
				if !continueOnError {
					return true
				}
			}
			for k, v := range outputs {
				tc = tc.WithBinding(k, v)
			}
		}
	}
	return false
}

func tryOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.Operation, cleaner cleaner.CleanerCollector) ([]operation, error) {
	var ops []operation
	if handler.Apply != nil {
		loaded, err := applyOperation(compilers, basePath, id+1, namespacer, cleaner, bindings, *handler.Apply)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Assert != nil {
		loaded, err := assertOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Assert)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(compilers, basePath, id+1, namespacer, *handler.Command))
	} else if handler.Create != nil {
		loaded, err := createOperation(compilers, basePath, id+1, namespacer, cleaner, bindings, *handler.Create)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(compilers, basePath, id+1, namespacer, *handler.Describe))
	} else if handler.Error != nil {
		loaded, err := errorOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Error)
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
		ops = append(ops, getOperation(compilers, basePath, id+1, namespacer, get))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(compilers, basePath, id+1, namespacer, *handler.Get))
	} else if handler.Patch != nil {
		loaded, err := patchOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Patch)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.PodLogs != nil {
		ops = append(ops, logsOperation(compilers, basePath, id+1, namespacer, *handler.PodLogs))
	} else if handler.Proxy != nil {
		ops = append(ops, proxyOperation(compilers, basePath, id+1, namespacer, *handler.Proxy))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(compilers, basePath, id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(compilers, id+1, *handler.Sleep))
	} else if handler.Update != nil {
		loaded, err := updateOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Update)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(compilers, basePath, id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func catchOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(compilers, basePath, id+1, namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(compilers, basePath, id+1, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(compilers, basePath, id+1, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(compilers, basePath, id+1, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(compilers, basePath, id+1, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(compilers, basePath, id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(compilers, id+1, *handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(compilers, basePath, id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func finallyOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]operation, error) {
	var ops []operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(compilers, basePath, id+1, namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(compilers, basePath, id+1, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(compilers, basePath, id+1, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(compilers, basePath, id+1, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, id+1, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(compilers, basePath, id+1, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(compilers, basePath, id+1, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(compilers, id+1, *handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(compilers, basePath, id+1, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func applyOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Apply) ([]operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:   basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
					timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Apply.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func assertOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Assert) ([]operation, error) {
	resources, err := fileRefOrCheck(context.TODO(), op.ActionCheckRef, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:   basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					templating: op.Template,
					timeouts:   &v1alpha1.Timeouts{Assert: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Assert.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func commandOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.Command) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				op := opcommand.New(
					tc.Compilers(),
					op,
					basePath,
					ns,
					config,
				)
				return op, ptr.To(tc.Timeouts().Exec.Duration), tc, nil
			}
		},
	)
}

func createOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Create) ([]operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:   basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
					timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Apply.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func deleteOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Delete) ([]operation, error) {
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
	resources, err := fileRefOrResource(context.TODO(), ref, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:            basePath,
					cluster:             op.Cluster,
					clusters:            op.Clusters,
					deletionPropagation: op.DeletionPropagationPolicy,
					templating:          op.Template,
					timeouts:            &v1alpha1.Timeouts{Delete: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Delete.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func describeOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.Describe) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData); err != nil {
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
					basePath,
					ns,
					config,
				)
				return op, ptr.To(tc.Timeouts().Exec.Duration), tc, nil
			}
		},
	)
}

func errorOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Error) ([]operation, error) {
	resources, err := fileRefOrCheck(context.TODO(), op.ActionCheckRef, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:   basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					templating: op.Template,
					timeouts:   &v1alpha1.Timeouts{Error: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Error.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func getOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.Get) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData); err != nil {
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
					basePath,
					ns,
					config,
				)
				return op, ptr.To(tc.Timeouts().Exec.Duration), tc, nil
			}
		},
	)
}

func logsOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.PodLogs) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData); err != nil {
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
					basePath,
					ns,
					config,
				)
				return op, ptr.To(tc.Timeouts().Exec.Duration), tc, nil
			}
		},
	)
}

func patchOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Patch) ([]operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:   basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
					timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Apply.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func proxyOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.Proxy) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData); err != nil {
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
					basePath,
					ns,
					config,
				)
				return op, ptr.To(tc.Timeouts().Exec.Duration), tc, nil
			}
		},
	)
}

func scriptOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.Script) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeScript,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
				return nil, nil, tc, err
			} else if config, _, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				op := opscript.New(
					tc.Compilers(),
					op,
					basePath,
					ns,
					config,
				)
				return op, ptr.To(tc.Timeouts().Exec.Duration), tc, nil
			}
		},
	)
}

func sleepOperation(_ compilers.Compilers, id int, op v1alpha1.Sleep) operation {
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeSleep,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			return opsleep.New(op), nil, tc, nil
		},
	)
}

func updateOperation(compilers compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Update) ([]operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
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
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := contextData{
					basePath:   basePath,
					cluster:    op.Cluster,
					clusters:   op.Clusters,
					dryRun:     op.DryRun,
					templating: op.Template,
					timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := setupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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
					return op, ptr.To(tc.Timeouts().Apply.Duration), tc, nil
				}
			},
		))
	}
	return ops, nil
}

func waitOperation(_ compilers.Compilers, basePath string, id int, namespacer namespacer.Namespacer, op v1alpha1.Wait) operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		OperationInfo{
			Id: id,
		},
		model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := contextData{
				basePath: basePath,
				cluster:  op.Cluster,
				clusters: op.Clusters,
				timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := setupContextAndBindings(tc, contextData); err != nil {
				return nil, nil, tc, err
			} else if config, client, err := tc.CurrentClusterClient(); err != nil {
				return nil, nil, tc, err
			} else {
				// make sure timeout is set to populate the command flag
				timeout := tc.Timeouts().Exec.Duration
				op.Timeout = &metav1.Duration{Duration: timeout}
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
					basePath,
					ns,
					config,
				)
				// shift operation timeout
				timeout += 30 * time.Second
				return op, &timeout, tc, nil
			}
		},
	)
}

func fileRefOrCheck(ctx context.Context, ref v1alpha1.ActionCheckRef, basePath string, compilers compilers.Compilers, bindings apis.Bindings) ([]unstructured.Unstructured, error) {
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
			return resource.Load(filepath.Join(basePath, ref), false)
		} else {
			return resource.LoadFromURI(url, false)
		}
	}
	return nil, errors.New("file or resource must be set")
}

func fileRefOrResource(ctx context.Context, ref v1alpha1.ActionResourceRef, basePath string, compilers compilers.Compilers, bindings apis.Bindings) ([]unstructured.Unstructured, error) {
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
			return resource.Load(filepath.Join(basePath, ref), true)
		} else {
			return resource.LoadFromURI(url, true)
		}
	}
	return nil, errors.New("file or resource must be set")
}

func prepareResource(resource unstructured.Unstructured, tc enginecontext.TestContext) error {
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

func getCleanerOrNil(cleaner cleaner.CleanerCollector, tc enginecontext.TestContext) cleaner.CleanerCollector {
	if tc.DryRun() {
		return nil
	}
	if tc.SkipDelete() {
		return nil
	}
	return cleaner
}
