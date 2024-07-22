package dryrun

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client = client.Client

type dryRunClient struct {
	inner Client
}

func (c *dryRunClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	return c.inner.Create(ctx, obj, append(opts, ctrlclient.DryRunAll)...)
}

func (c *dryRunClient) Update(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
	return c.inner.Update(ctx, obj, append(opts, ctrlclient.DryRunAll)...)
}

func (c *dryRunClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	return c.inner.Delete(ctx, obj, append(opts, ctrlclient.DryRunAll)...)
}

func (c *dryRunClient) Get(ctx context.Context, key types.NamespacedName, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	return c.inner.Get(ctx, key, obj, opts...)
}

func (c *dryRunClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.inner.IsObjectNamespaced(obj)
}

func (c *dryRunClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
	return c.inner.List(ctx, list, opts...)
}

func (c *dryRunClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
	return c.inner.Patch(ctx, obj, patch, append(opts, ctrlclient.DryRunAll)...)
}

func (c *dryRunClient) RESTMapper() meta.RESTMapper {
	return c.inner.RESTMapper()
}

func New(inner Client) Client {
	return &dryRunClient{inner: inner}
}
