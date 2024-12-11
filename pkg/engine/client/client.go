package client

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func New(inner client.Client) client.Client {
	return &runnerClient{
		inner: inner,
	}
}

type runnerClient struct {
	inner client.Client
}

func (c *runnerClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) (_err error) {
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

func (c *runnerClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.ok(ctx, logging.Patch, obj)
		} else {
			c.error(ctx, logging.Update, obj, _err)
		}
	}()
	return c.inner.Update(ctx, obj, opts...)
}

func (c *runnerClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) (_err error) {
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

func (c *runnerClient) Get(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
	return c.inner.Get(ctx, key, obj, opts...)
}

func (c *runnerClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) (_err error) {
	return c.inner.List(ctx, list, opts...)
}

func (c *runnerClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) (_err error) {
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

func (c *runnerClient) RESTMapper() meta.RESTMapper {
	return c.inner.RESTMapper()
}

func (c *runnerClient) ok(ctx context.Context, op logging.Operation, obj client.Object) {
	logging.Log(ctx, op, logging.OkStatus, obj, color.BoldGreen)
}

func (c *runnerClient) error(ctx context.Context, op logging.Operation, obj client.Object, err error) {
	logging.Log(ctx, op, logging.OkStatus, obj, color.BoldYellow, logging.ErrSection(err))
}
