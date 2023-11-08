package operations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func operationCreate(ctx context.Context, obj ctrlclient.Object, c client.Client, shouldFail bool, cleaner cleanup.Cleaner) (_err error) {
	const operation = "CREATE"
	logger := logging.FromContext(ctx).WithResource(obj)
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
		actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
		err := c.Get(ctx, client.ObjectKey(obj), &actual)
		if err == nil {
			return false, errors.New("the resource already exists in the cluster")
		} else if kerrors.IsNotFound(err) {
			if err := c.Create(ctx, obj); err != nil {
				if shouldFail {
					return true, nil
				}
				return false, err
			} else {
				if obj.GetGeneration() != 0 {
					statusBytes, err := json.Marshal(obj)
					if err != nil {
						return false, err
					}
					if err := c.Status().Patch(ctx, obj, ctrlclient.RawPatch(types.MergePatchType, statusBytes)); err != nil {
						if shouldFail {
							return true, nil
						}
						return false, err
					} else if shouldFail {
						return false, errors.New("an error was expected but didn't happen")
					}
				}

				if cleaner != nil {
					cleaner(obj, c)
				}
				if shouldFail {
					return false, errors.New("an error was expected but didn't happen")
				}
			}
		} else {
			return false, err
		}
		return true, nil
	})
}
