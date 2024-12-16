package operations

import (
	"context"
	"errors"
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
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
)

func TryOperation(
	tc enginecontext.TestContext,
	basePath string,
	namespacer namespacer.Namespacer,
	handler v1alpha1.Operation,
	cleaner cleaner.CleanerCollector,
) (model.OperationType, []Operation, error) {
	if handler.Apply != nil {
		loaded, err := applyOperation(tc.Compilers(), basePath, namespacer, cleaner, tc.Bindings(), *handler.Apply)
		return model.OperationTypeApply, loaded, err
	} else if handler.Assert != nil {
		loaded, err := assertOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Assert)
		return model.OperationTypeAssert, loaded, err
	} else if handler.Command != nil {
		return model.OperationTypeCommand, []Operation{commandOperation(tc.Compilers(), basePath, namespacer, *handler.Command)}, nil
	} else if handler.Create != nil {
		loaded, err := createOperation(tc.Compilers(), basePath, namespacer, cleaner, tc.Bindings(), *handler.Create)
		return model.OperationTypeCreate, loaded, err
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Delete)
		return model.OperationTypeDelete, loaded, err
	} else if handler.Describe != nil {
		return model.OperationTypeCommand, []Operation{describeOperation(tc.Compilers(), basePath, namespacer, *handler.Describe)}, nil
	} else if handler.Error != nil {
		loaded, err := errorOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Error)
		return model.OperationTypeError, loaded, err
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
		return model.OperationTypeCommand, []Operation{getOperation(tc.Compilers(), basePath, namespacer, get)}, nil
	} else if handler.Get != nil {
		return model.OperationTypeCommand, []Operation{getOperation(tc.Compilers(), basePath, namespacer, *handler.Get)}, nil
	} else if handler.Patch != nil {
		loaded, err := patchOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Patch)
		return model.OperationTypePatch, loaded, err
	} else if handler.PodLogs != nil {
		return model.OperationTypeCommand, []Operation{logsOperation(tc.Compilers(), basePath, namespacer, *handler.PodLogs)}, nil
	} else if handler.Proxy != nil {
		return model.OperationTypeCommand, []Operation{proxyOperation(tc.Compilers(), basePath, namespacer, *handler.Proxy)}, nil
	} else if handler.Script != nil {
		return model.OperationTypeScript, []Operation{scriptOperation(tc.Compilers(), basePath, namespacer, *handler.Script)}, nil
	} else if handler.Sleep != nil {
		return model.OperationTypeSleep, []Operation{sleepOperation(tc.Compilers(), *handler.Sleep)}, nil
	} else if handler.Update != nil {
		loaded, err := updateOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Update)
		return model.OperationTypeUpdate, loaded, err
	} else if handler.Wait != nil {
		return model.OperationTypeCommand, []Operation{waitOperation(tc.Compilers(), basePath, namespacer, *handler.Wait)}, nil
	} else {
		return "", nil, errors.New("no operation found")
	}
}

func CatchOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]Operation, error) {
	var ops []Operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(compilers, basePath, namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(compilers, basePath, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(compilers, basePath, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(compilers, basePath, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(compilers, basePath, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(compilers, basePath, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(compilers, *handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(compilers, basePath, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func FinallyOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]Operation, error) {
	var ops []Operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(compilers, basePath, namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(compilers, basePath, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(compilers, basePath, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(compilers, basePath, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(compilers, basePath, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(compilers, basePath, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(compilers, *handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(compilers, basePath, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func applyOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Apply) ([]Operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypeApply,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:   basePath,
					Cluster:    op.Cluster,
					Clusters:   op.Clusters,
					DryRun:     op.DryRun,
					Templating: op.Template,
					Timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func assertOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Assert) ([]Operation, error) {
	resources, err := fileRefOrCheck(context.TODO(), op.ActionCheckRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypeAssert,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:   basePath,
					Cluster:    op.Cluster,
					Clusters:   op.Clusters,
					Templating: op.Template,
					Timeouts:   &v1alpha1.Timeouts{Assert: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func commandOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.Command) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func createOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Create) ([]Operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypeCreate,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:   basePath,
					Cluster:    op.Cluster,
					Clusters:   op.Clusters,
					DryRun:     op.DryRun,
					Templating: op.Template,
					Timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func deleteOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Delete) ([]Operation, error) {
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
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypeDelete,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:            basePath,
					Cluster:             op.Cluster,
					Clusters:            op.Clusters,
					DeletionPropagation: op.DeletionPropagationPolicy,
					Templating:          op.Template,
					Timeouts:            &v1alpha1.Timeouts{Delete: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func describeOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.Describe) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
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

func errorOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Error) ([]Operation, error) {
	resources, err := fileRefOrCheck(context.TODO(), op.ActionCheckRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypeError,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:   basePath,
					Cluster:    op.Cluster,
					Clusters:   op.Clusters,
					Templating: op.Template,
					Timeouts:   &v1alpha1.Timeouts{Error: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func getOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.Get) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
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

func logsOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.PodLogs) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
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

func patchOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Patch) ([]Operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypePatch,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:   basePath,
					Cluster:    op.Cluster,
					Clusters:   op.Clusters,
					DryRun:     op.DryRun,
					Templating: op.Template,
					Timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func proxyOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.Proxy) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
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

func scriptOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.Script) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeScript,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func sleepOperation(_ compilers.Compilers, op v1alpha1.Sleep) Operation {
	return newOperation(
		// model.OperationTypeSleep,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			return opsleep.New(op), nil, tc, nil
		},
	)
}

func updateOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Update) ([]Operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, newOperation(
			// model.OperationTypeUpdate,
			func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
				contextData := enginecontext.ContextData{
					BasePath:   basePath,
					Cluster:    op.Cluster,
					Clusters:   op.Clusters,
					DryRun:     op.DryRun,
					Templating: op.Template,
					Timeouts:   &v1alpha1.Timeouts{Apply: op.Timeout},
				}
				if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, op.Bindings...); err != nil {
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

func waitOperation(_ compilers.Compilers, basePath string, namespacer namespacer.Namespacer, op v1alpha1.Wait) Operation {
	ns := ""
	if namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	return newOperation(
		// model.OperationTypeCommand,
		func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
			contextData := enginecontext.ContextData{
				BasePath: basePath,
				Cluster:  op.Cluster,
				Clusters: op.Clusters,
				Timeouts: &v1alpha1.Timeouts{Exec: op.Timeout},
			}
			if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
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
