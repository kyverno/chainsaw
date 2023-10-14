package client

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/utils/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Assert(ctx context.Context, client Client, expected ctrlclient.Object) error {

	// Resource do have a namespace if they are namespaced

	// We assume here that the gvk kind already exist in the cluster.

	// Step 1: Get the object from the cluster
	// Step 2 We need to make sure that obj is cluster scoped or namespaced
	// Step 3: Compare the object with the expected object
	// Step 4: Return true if the objects are equal, false otherwise

	// gvk := expected.GetObjectKind().GroupVersionKind()

	actualObj := unstructured.Unstructured{}
	expectedObj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(expected)

	name, ns, err := getNameandNamespace(expected)
	if err != nil {
		return err
	}

	err = client.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: ns}, &actualObj)
	if err != nil {
		return err
	}

	kubernetes.IsSubset(actualObj, expectedObj)
	return nil
}

func getNameandNamespace(obj ctrlclient.Object) (string, string, error) {

	if obj.GetName() == "" {
		return "", "", fmt.Errorf("object does not have a name")
	}

	return obj.GetName(), obj.GetNamespace(), nil

}
