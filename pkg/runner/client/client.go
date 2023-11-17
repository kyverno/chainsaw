package client

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func New(inner client.Client) client.Client {
	return &runnerClient{
		inner: inner,
	}
}

type runnerClient struct {
	inner client.Client
}

func (c *runnerClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.ok(ctx, logging.Create, obj)
		} else {
			c.error(ctx, logging.Create, obj, _err)
		}
	}()
	err := c.inner.Create(ctx, obj, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (c *runnerClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.ok(ctx, logging.Delete, obj)
		} else {
			c.error(ctx, logging.Delete, obj, _err)
		}
	}()
	return c.inner.Delete(ctx, obj, opts...)
}

func (c *runnerClient) Get(ctx context.Context, key types.NamespacedName, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	return c.inner.Get(ctx, key, obj, opts...)
}

func (c *runnerClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) (_err error) {
	return c.inner.List(ctx, list, opts...)
}

func (c *runnerClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.ok(ctx, logging.Patch, obj)
		} else {
			c.error(ctx, logging.Patch, obj, _err)
		}
	}()
	return c.inner.Patch(ctx, obj, patch, opts...)
}

func (c *runnerClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.inner.IsObjectNamespaced(obj)
}

func (c *runnerClient) ok(ctx context.Context, op logging.Operation, obj ctrlclient.Object) {
	logger := logging.FromContext(ctx)
	if logger != nil {
		logger.WithResource(obj).Log(op, logging.OkStatus, color.BoldGreen)
	}
}

func (c *runnerClient) error(ctx context.Context, op logging.Operation, obj ctrlclient.Object, err error) {
	logger := logging.FromContext(ctx)
	if logger != nil {
		logger.WithResource(obj).Log(op, logging.ErrorStatus, color.BoldYellow, logging.ErrSection(err))
	}
}
