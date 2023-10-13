package client

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/utils/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Assert(ctx context.Context, client Client, expected ctrlclient.Object, namespace string) (bool, error) {

	// Step 1: Get the object from the cluster
	// Step 2 We need to make sure that obj is cluster scoped or namespaced
	// Step 3: Compare the object with the expected object
	// Step 4: Return true if the objects are equal, false otherwise

	// Try to check if the object is namespaced or cluster scoped and get the name and namespace
	name, ns, err := kubernetes.Namespaced(expected, namespace)
	if err != nil {
		return false, err
	}

	gvk := expected.GetObjectKind().GroupVersionKind()

	actualObj := unstructured.Unstructured{}

	err = client.Get(ctx, ctrlclient.ObjectKey{Name: name, Namespace: ns}, &actualObj)
	if err != nil {
		return false, err
	}

	expectedObj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(expected)

	kubernetes.IsSubset(actualObj, expectedObj)
	return true, nil
}
