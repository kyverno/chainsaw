package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Delete(ctx context.Context, expected unstructured.Unstructured, c client.Client) error {
	candidates, err := read(ctx, expected, c)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	for _, candidate := range candidates {
		err := c.Delete(ctx, &candidate)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	for _, candidate := range candidates {
		if err := wait.PollUntilContextCancel(ctx, interval, true, func(ctx context.Context) (bool, error) {
			var actual unstructured.Unstructured
			actual.SetGroupVersionKind(candidate.GetObjectKind().GroupVersionKind())
			err := c.Get(ctx, client.ObjectKey(&candidate), &actual)
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
	return nil
}
