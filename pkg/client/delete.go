package client

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func BlockingDelete(ctx context.Context, client Client, obj ctrlclient.Object) error {
	err := client.Delete(ctx, obj)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	for {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
		err := client.Get(ctx, ObjectKey(obj), &actual)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}
	}
	// TODO: context timeout
	// return nil
}
