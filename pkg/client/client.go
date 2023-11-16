package client

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client interface {
	// Get retrieves an obj for the given object key from the Kubernetes Cluster.
	// obj must be a struct pointer so that obj can be updated with the response
	// returned by the Server.
	Get(ctx context.Context, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error

	// List retrieves list of objects for a given namespace and list options. On a
	// successful call, Items field in the list will be populated with the
	// result returned from the server.
	List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error

	// Create saves the object obj in the Kubernetes cluster. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error

	// Delete deletes the given obj from Kubernetes cluster.
	Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error

	// Patch patches the given obj in the Kubernetes cluster. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error

	// IsObjectNamespaced returns true if the GroupVersionKind of the object is namespaced.
	IsObjectNamespaced(obj runtime.Object) (bool, error)

	// Status update the status subresource for the given object.
	Status() ctrlclient.StatusWriter
}

func New(cfg *rest.Config) (Client, error) {
	var opts ctrlclient.Options
	return ctrlclient.New(cfg, opts)
}
