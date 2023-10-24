package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Delete(ctx context.Context, logger logging.Logger, expected ctrlclient.Object, c client.Client) (_err error) {
	attempts := 0
	defer func() {
		logging.ResourceOp(logger, "DELETE", client.ObjectKey(expected), expected, attempts, _err)
	}()
	candidates, _err := read(ctx, expected, c)
	if _err != nil {
		if errors.IsNotFound(_err) {
			return nil
		}
		return _err
	}
	for i := range candidates {
		if candidates[i].GetName() != "kube-root-ca.crt" {
			err := c.Delete(ctx, &candidates[i])
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}
	}
	gvk := expected.GetObjectKind().GroupVersionKind()
	for i := range candidates {
		if candidates[i].GetName() != "kube-root-ca.crt" {
			if err := wait.PollUntilContextCancel(ctx, interval, true, func(ctx context.Context) (bool, error) {
				attempts++
				var actual unstructured.Unstructured
				actual.SetGroupVersionKind(gvk)
				err := c.Get(ctx, client.ObjectKey(&candidates[i]), &actual)
				if err != nil {
					if errors.IsNotFound(err) {
						return true, nil
					}
					return false, err
				}
				return false, nil
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
