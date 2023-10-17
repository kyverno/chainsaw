package client

import (
	"context"
	"fmt"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
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

func DeleteResource(ctx context.Context, client Client, delete *unstructured.Unstructured) error {
	err := client.Delete(context.TODO(), delete)
	if err != nil && !kerror.IsNotFound(err) {
		return err
	}
	// Wait for resources to be deleted.
	return wait.PollUntilContextTimeout(ctx, 100*time.Millisecond, 1*time.Second, true, wait.ConditionWithContextFunc(func(ctx context.Context) (bool, error) {
		actual := &unstructured.Unstructured{}
		actual.SetGroupVersionKind(delete.GetObjectKind().GroupVersionKind())
		err = client.Get(context.TODO(), types.NamespacedName{Name: delete.GetName(), Namespace: delete.GetNamespace()}, actual)
		if err == nil || !kerror.IsNotFound(err) {
			return false, err
		}
		return true, nil
	}))
}

func NewResource(apiVersion, kind, name, namespace string) *unstructured.Unstructured {
	meta := map[string]interface{}{}

	if name != "" {
		meta["name"] = name
	}
	if namespace != "" {
		meta["namespace"] = namespace
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata":   meta,
		},
	}
}

func ListResourcesToDelete(c Client, delete v1alpha1.ObjectReference) (*unstructured.UnstructuredList, error) {
	u := &unstructured.UnstructuredList{}
	u.SetGroupVersionKind(delete.GetObjectKind().GroupVersionKind())

	listOptions := []ctrlclient.ListOption{}
	if delete.Labels != nil {
		listOptions = append(listOptions, ctrlclient.MatchingLabels(delete.Labels))
	}

	if delete.Namespace != "" {
		listOptions = append(listOptions, ctrlclient.InNamespace(delete.Namespace))
	}

	err := c.List(context.TODO(), u, listOptions...)
	if err != nil {
		return nil, fmt.Errorf("listing matching resources: %v", err)
	}
	return u, nil
}
