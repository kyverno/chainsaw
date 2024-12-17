package operations

import (
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
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
		return model.OperationTypeCommand, []Operation{commandOperation(basePath, namespacer, *handler.Command)}, nil
	} else if handler.Create != nil {
		loaded, err := createOperation(tc.Compilers(), basePath, namespacer, cleaner, tc.Bindings(), *handler.Create)
		return model.OperationTypeCreate, loaded, err
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Delete)
		return model.OperationTypeDelete, loaded, err
	} else if handler.Describe != nil {
		return model.OperationTypeCommand, []Operation{describeOperation(basePath, namespacer, *handler.Describe)}, nil
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
		return model.OperationTypeCommand, []Operation{getOperation(basePath, namespacer, get)}, nil
	} else if handler.Get != nil {
		return model.OperationTypeCommand, []Operation{getOperation(basePath, namespacer, *handler.Get)}, nil
	} else if handler.Patch != nil {
		loaded, err := patchOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Patch)
		return model.OperationTypePatch, loaded, err
	} else if handler.PodLogs != nil {
		return model.OperationTypeCommand, []Operation{logsOperation(basePath, namespacer, *handler.PodLogs)}, nil
	} else if handler.Proxy != nil {
		return model.OperationTypeCommand, []Operation{proxyOperation(basePath, namespacer, *handler.Proxy)}, nil
	} else if handler.Script != nil {
		return model.OperationTypeScript, []Operation{scriptOperation(basePath, namespacer, *handler.Script)}, nil
	} else if handler.Sleep != nil {
		return model.OperationTypeSleep, []Operation{sleepOperation(*handler.Sleep)}, nil
	} else if handler.Update != nil {
		loaded, err := updateOperation(tc.Compilers(), basePath, namespacer, tc.Bindings(), *handler.Update)
		return model.OperationTypeUpdate, loaded, err
	} else if handler.Wait != nil {
		return model.OperationTypeCommand, []Operation{waitOperation(basePath, namespacer, *handler.Wait)}, nil
	} else {
		return "", nil, errors.New("no operation found")
	}
}

func CatchOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]Operation, error) {
	var ops []Operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(basePath, namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(basePath, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(basePath, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(basePath, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(basePath, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(basePath, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(*handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(basePath, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}

func FinallyOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, handler v1alpha1.CatchFinally) ([]Operation, error) {
	var ops []Operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(basePath, namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(basePath, namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(basePath, namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(basePath, namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(compilers, basePath, namespacer, bindings, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(basePath, namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(basePath, namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(*handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(basePath, namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}
