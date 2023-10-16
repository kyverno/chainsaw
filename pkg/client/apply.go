package client

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // package needed for auth providers like GCP
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateOrUpdate(ctx context.Context, client Client, obj ctrlclient.Object) error {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := client.Get(ctx, ObjectKey(obj), &actual)
	if err == nil {
		bytes, err := PatchObject(&actual, obj)
		if err != nil {
			return err
		}
		return client.Patch(ctx, &actual, ctrlclient.RawPatch(types.MergePatchType, bytes))
	} else if errors.IsNotFound(err) {
		return client.Create(ctx, obj)
	}
	// TODO: context timeout
	return err
}
