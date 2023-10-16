package client

import (
	"context"
	"errors"
	"strings"

	"github.com/kyverno/chainsaw/pkg/utils/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dClient "k8s.io/client-go/discovery"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Assert(ctx context.Context, expected ctrlclient.Object, client Client, dClient dClient.DiscoveryInterface) error {
	// Resource do have a namespace if they are namespaced
	// Step 1: Check if the resource exists in the cluster
	// Step 2 Get the actual object from the cluster
	// Step 3: Compare the actual object with the expected object
	// Check if the resource exists in the cluster
	gvk := expected.GetObjectKind().GroupVersionKind()
	if err := checkAPIResource(dClient, gvk); err != nil {
		return err
	}

	actualObj := unstructured.Unstructured{}
	expectedObj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(expected)

	if err := client.Get(ctx, ctrlclient.ObjectKey{Name: expected.GetName(), Namespace: expected.GetNamespace()}, &actualObj); err != nil {
		return err
	}

	kubernetes.IsSubset(actualObj, expectedObj)
	return nil
}

// getAPIResource returns the APIResource object for a specific GroupVersionKind.
func checkAPIResource(dClient dClient.DiscoveryInterface, gvk schema.GroupVersionKind) error {
	resourceTypes, err := dClient.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return err
	}

	for _, resource := range resourceTypes.APIResources {
		if !strings.EqualFold(resource.Kind, gvk.Kind) {
			continue
		}
		return nil
	}

	return errors.New("resource type not found")
}
