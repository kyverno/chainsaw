package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Apply(ctx context.Context, logger logging.Logger, obj ctrlclient.Object, c client.Client, cleanup CleanupFunc) error {
	logging.ResourceOp(logger, "APPLY", client.ObjectKey(obj), obj)
	return wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
		err := c.Get(ctx, client.ObjectKey(obj), &actual)
		if err == nil {
			bytes, err := client.PatchObject(&actual, obj)
			if err != nil {
				return false, err
			}
			if err := c.Patch(ctx, &actual, ctrlclient.RawPatch(types.MergePatchType, bytes)); err != nil {
				return false, err
			}
		} else if errors.IsNotFound(err) {
			if err := c.Create(ctx, obj); err != nil {
				return false, err
			}
			if cleanup != nil {
				cleanup(obj, c)
			}
		}
		return true, nil
	})
}
