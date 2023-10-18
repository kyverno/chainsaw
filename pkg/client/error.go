package client

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/utils/kubernetes"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Error(ctx context.Context, expected unstructured.Unstructured, client Client) error {
	// Resource do have a namespace if they are namespaced
	// Step 1: Check if the resource exists in the cluster using metadata
	// Step 2: Get the actual object from the cluster return nil if not found in the cluster
	// Step 3: Compare the actual object with the expected object return err if equal

	var actuals []unstructured.Unstructured

	name := expected.GetName()
	namespace := expected.GetNamespace()
	gvk := expected.GetObjectKind().GroupVersionKind()

	if name != "" {
		actual := unstructured.Unstructured{}
		actual.SetGroupVersionKind(gvk)

		if err := client.Get(ctx, ctrlclient.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}, &actual); err != nil {
			if kerror.IsNotFound(err) {
				return nil
			}
			return err
		}

		actuals = []unstructured.Unstructured{actual}
	} else {
		var err error
		actuals, err = list(client, gvk, namespace)
		if err != nil {
			return err
		}
	}

	var unexpectedObjects []unstructured.Unstructured
	for _, actual := range actuals {
		if err := kubernetes.IsSubset(expected.UnstructuredContent(), actual.UnstructuredContent()); err == nil {
			unexpectedObjects = append(unexpectedObjects, actual)
		}
	}

	if len(unexpectedObjects) == 0 {
		return nil
	}
	if len(unexpectedObjects) == 1 {
		return fmt.Errorf("resource %s %s matched error assertion", unexpectedObjects[0].GroupVersionKind(), unexpectedObjects[0].GetName())
	}
	return fmt.Errorf("resource %s %s (and %d other resources) matched error assertion", unexpectedObjects[0].GroupVersionKind(), unexpectedObjects[0].GetName(), len(unexpectedObjects)-1)
}

func list(cl Client, gvk schema.GroupVersionKind, namespace string) ([]unstructured.Unstructured, error) {
	list := unstructured.UnstructuredList{}
	list.SetGroupVersionKind(gvk)

	listOptions := []ctrlclient.ListOption{}
	if namespace != "" {
		listOptions = append(listOptions, ctrlclient.InNamespace(namespace))
	}

	if err := cl.List(context.TODO(), &list, listOptions...); err != nil {
		return []unstructured.Unstructured{}, err
	}

	return list.Items, nil
}
