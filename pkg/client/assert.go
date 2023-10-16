package client

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/utils/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Assert(ctx context.Context, expected unstructured.Unstructured, client Client) error {
	// TODO: we should retry until context timeout
	actual := unstructured.Unstructured{}
	actual.SetGroupVersionKind(expected.GetObjectKind().GroupVersionKind())
	if err := client.Get(ctx, ctrlclient.ObjectKey{Name: expected.GetName(), Namespace: expected.GetNamespace()}, &actual); err != nil {
		return err
	}
	if kubernetes.IsSubset(actual, expected) {
		return nil
	}
	return fmt.Errorf("expected object %v is not subset of actual object %v", expected, actual)
}
