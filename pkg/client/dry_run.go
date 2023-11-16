package client

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type dryRunClient struct {
	inner Client
}

func (c *dryRunClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return c.inner.Create(ctx, obj, append(opts, client.DryRunAll)...)
}

func (c *dryRunClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return c.inner.Delete(ctx, obj, append(opts, client.DryRunAll)...)
}

func (c *dryRunClient) Get(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
	return c.inner.Get(ctx, key, obj, opts...)
}

func (c *dryRunClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.inner.IsObjectNamespaced(obj)
}

func (c *dryRunClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return c.inner.List(ctx, list, opts...)
}

func (c *dryRunClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return c.inner.Patch(ctx, obj, patch, append(opts, client.DryRunAll)...)
}

// Don't follow dry-run for status updates
func (c *dryRunClient) Status() ctrlclient.StatusWriter {
	return c.inner.Status()
}

func DryRun(inner Client) Client {
	return &dryRunClient{inner: inner}
}
