package operations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplyOperation struct {
	BaseOperation
}

func (a *ApplyOperation) Name() string {
	return "APPLY"
}

func (a *ApplyOperation) Exec(ctx context.Context) (_err error) {
	const operation = "APPLY"
	logger := logging.FromContext(ctx).WithResource(a.obj)
	logger.Log(operation, color.BoldFgCyan, "RUNNING...")
	defer func() {
		if _err == nil {
			logger.Log(operation, color.BoldGreen, "DONE")
		} else {
			logger.Log(operation, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()

	return wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(a.obj.GetObjectKind().GroupVersionKind())
		err := a.client.Get(ctx, client.ObjectKey(a.obj), &actual)
		if err == nil {
			patchOptions := []ctrlclient.PatchOption{}
			if a.dryRun {
				patchOptions = append(patchOptions, ctrlclient.DryRunAll)
			}

			patched, err := client.PatchObject(&actual, a.obj)
			if err != nil {
				return false, err
			}
			bytes, err := json.Marshal(patched)
			if err != nil {
				return false, err
			}

			if err := a.client.Patch(ctx, &actual, ctrlclient.RawPatch(types.MergePatchType, bytes), patchOptions...); err != nil {
				if a.shouldFail {
					return true, nil
				}
				return false, err
			} else if a.shouldFail {
				return false, errors.New("an error was expected but didn't happen")
			}
			if a.dryRun {
				logger.Log(operation, color.BoldYellow, "DRY RUN: Resource patch simulated")
			}
		} else if kerrors.IsNotFound(err) {
			var createOptions []ctrlclient.CreateOption
			if a.dryRun {
				createOptions = append(createOptions, ctrlclient.DryRunAll)
			}
			if err := a.client.Create(ctx, a.obj, createOptions...); err != nil {
				if a.shouldFail {
					return true, nil
				}
				return false, err
			} else {
				if a.cleaner != nil && !a.dryRun {
					a.cleaner(a.obj, a.client)
				}
				if a.shouldFail {
					return false, errors.New("an error was expected but didn't happen")
				}
				if a.dryRun {
					logger.Log(operation, color.BoldYellow, "DRY RUN: Resource creation simulated")
				}
			}
		} else {
			return false, err
		}
		return true, nil
	})
}
