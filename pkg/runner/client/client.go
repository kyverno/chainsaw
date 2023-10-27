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

func (c *runnerClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.log(color.BoldGreen.Sprint("CREATE"), obj, "OK")
		} else {
			c.log(color.BoldRed.Sprint("CREATE"), obj, "ERROR", _err)
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
			c.log(color.BoldGreen.Sprint("DELETE"), obj, "OK")
		} else {
			c.log(color.BoldRed.Sprint("DELETE"), obj, "ERROR", _err)
		}
	}()
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

func (c *runnerClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) (_err error) {
	return c.inner.List(ctx, list, opts...)
}

func (c *runnerClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.log(color.BoldGreen.Sprint("PATCH"), obj, "OK")
		} else {
			c.log(color.BoldRed.Sprint("PATCH"), obj, "ERROR", _err)
		}
	}()
	return c.inner.Patch(ctx, obj, patch, opts...)
}

func (c *runnerClient) Update(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) (_err error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	defer func() {
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		if _err == nil {
			c.log(color.BoldGreen.Sprint("UPDATE"), obj, "OK")
		} else {
			c.log(color.BoldRed.Sprint("UPDATE"), obj, "ERROR", _err)
		}
	}()
	return c.inner.Update(ctx, obj, opts...)
}

func (c *runnerClient) log(op string, obj ctrlclient.Object, args ...interface{}) {
	c.logger.WithResource(obj).Log(op, args...)
}
