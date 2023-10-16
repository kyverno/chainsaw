package client

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/utils/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Assert(ctx context.Context, expected ctrlclient.Object, client Client) error {
	// Resource do have a namespace if they are namespaced
	// Step 1: Check if the resource exists in the cluster
	// Step 2 Get the actual object from the cluster
	// Step 3: Compare the actual object with the expected object
	// Check if the resource exists in the cluster
	actualObj := unstructured.Unstructured{}
	expectedObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(expected)
	if err != nil {
		return err
	}

	if err := client.Get(ctx, ctrlclient.ObjectKey{Name: expected.GetName(), Namespace: expected.GetNamespace()}, &actualObj); err != nil {
		return err
	}

	if kubernetes.IsSubset(actualObj, expectedObj) {
		return nil
	}
	return fmt.Errorf("expected object %v is not subset of actual object %v", expectedObj, actualObj)
}
