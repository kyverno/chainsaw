package operations

import (
	"context"
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

func TryOperation(
	ctx context.Context,
	tc enginecontext.TestContext,
	namespacer namespacer.Namespacer,
	handler v1alpha1.Operation,
	cleaner cleaner.CleanerCollector,
) (model.OperationType, []Operation, error) {
	if handler.Apply != nil {
		loaded, err := applyOperation(ctx, tc, namespacer, cleaner, *handler.Apply)
		return model.OperationTypeApply, loaded, err
	} else if handler.Assert != nil {
		loaded, err := assertOperation(ctx, tc, namespacer, *handler.Assert)
		return model.OperationTypeAssert, loaded, err
	} else if handler.Command != nil {
		return model.OperationTypeCommand, []Operation{commandOperation(namespacer, *handler.Command)}, nil
	} else if handler.Create != nil {
		loaded, err := createOperation(ctx, tc, namespacer, cleaner, *handler.Create)
		return model.OperationTypeCreate, loaded, err
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(ctx, tc, namespacer, *handler.Delete)
		return model.OperationTypeDelete, loaded, err
	} else if handler.Describe != nil {
		return model.OperationTypeCommand, []Operation{describeOperation(namespacer, *handler.Describe)}, nil
	} else if handler.Error != nil {
		loaded, err := errorOperation(ctx, tc, namespacer, *handler.Error)
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
		return model.OperationTypeCommand, []Operation{getOperation(namespacer, get)}, nil
	} else if handler.Get != nil {
		return model.OperationTypeCommand, []Operation{getOperation(namespacer, *handler.Get)}, nil
	} else if handler.Patch != nil {
		loaded, err := patchOperation(ctx, tc, namespacer, *handler.Patch)
		return model.OperationTypePatch, loaded, err
	} else if handler.PodLogs != nil {
		return model.OperationTypeCommand, []Operation{logsOperation(namespacer, *handler.PodLogs)}, nil
	} else if handler.Proxy != nil {
		return model.OperationTypeCommand, []Operation{proxyOperation(namespacer, *handler.Proxy)}, nil
	} else if handler.Script != nil {
		return model.OperationTypeScript, []Operation{scriptOperation(namespacer, *handler.Script)}, nil
	} else if handler.Sleep != nil {
		return model.OperationTypeSleep, []Operation{sleepOperation(*handler.Sleep)}, nil
	} else if handler.Update != nil {
		loaded, err := updateOperation(ctx, tc, namespacer, *handler.Update)
		return model.OperationTypeUpdate, loaded, err
	} else if handler.Wait != nil {
		return model.OperationTypeCommand, []Operation{waitOperation(namespacer, *handler.Wait)}, nil
	} else {
		return "", nil, errors.New("no operation found")
	}
}

func CatchOperation(
	ctx context.Context,
	tc enginecontext.TestContext,
	namespacer namespacer.Namespacer,
	handler v1alpha1.CatchFinally,
) ([]Operation, error) {
	var ops []Operation
	if handler.PodLogs != nil {
		ops = append(ops, logsOperation(namespacer, *handler.PodLogs))
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
		ops = append(ops, getOperation(namespacer, get))
	} else if handler.Describe != nil {
		ops = append(ops, describeOperation(namespacer, *handler.Describe))
	} else if handler.Get != nil {
		ops = append(ops, getOperation(namespacer, *handler.Get))
	} else if handler.Delete != nil {
		loaded, err := deleteOperation(ctx, tc, namespacer, *handler.Delete)
		if err != nil {
			return nil, err
		}
		ops = append(ops, loaded...)
	} else if handler.Command != nil {
		ops = append(ops, commandOperation(namespacer, *handler.Command))
	} else if handler.Script != nil {
		ops = append(ops, scriptOperation(namespacer, *handler.Script))
	} else if handler.Sleep != nil {
		ops = append(ops, sleepOperation(*handler.Sleep))
	} else if handler.Wait != nil {
		ops = append(ops, waitOperation(namespacer, *handler.Wait))
	} else {
		return nil, errors.New("no operation found")
	}
	return ops, nil
}