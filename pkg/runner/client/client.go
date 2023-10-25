package client

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func New(logger logging.Logger, inner client.Client) client.Client {
	return &runnerClient{
		logger: logger,
		inner:  inner,
	}
}

type runnerClient struct {
	logger logging.Logger
	inner  client.Client
}

func (c *runnerClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	c.log("create", client.ObjectKey(obj), obj)
	err := c.inner.Create(ctx, obj, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (c *runnerClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	return c.inner.Delete(ctx, obj, opts...)
}

func (c *runnerClient) DeleteAllOf(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteAllOfOption) error {
	return c.inner.DeleteAllOf(ctx, obj, opts...)
}

func (c *runnerClient) Get(ctx context.Context, key types.NamespacedName, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	return c.inner.Get(ctx, key, obj, opts...)
}

func (c *runnerClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.inner.IsObjectNamespaced(obj)
}

func (c *runnerClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
	return c.inner.List(ctx, list, opts...)
}

func (c *runnerClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
	c.log("patch", client.ObjectKey(obj), obj)
	return c.inner.Patch(ctx, obj, patch, opts...)
}

func (c *runnerClient) Update(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
	c.log("update", client.ObjectKey(obj), obj)
	return c.inner.Update(ctx, obj, opts...)
}

func (c *runnerClient) log(op string, key ctrlclient.ObjectKey, obj ctrlclient.Object) {
	// logging.ResourceOp(c.logger, op, key, obj)
}
