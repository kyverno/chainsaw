package client

import (
	"context"
	"fmt"

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
			c.log(ctx, logging.Create, obj, color.BoldGreen, "OK")
		} else {
			c.log(ctx, logging.Create, obj, color.BoldYellow, fmt.Sprintf("ERROR\n%s", _err))
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
			c.log(ctx, logging.Delete, obj, color.BoldGreen, "OK")
		} else {
			c.log(ctx, logging.Delete, obj, color.BoldYellow, fmt.Sprintf("ERROR\n%s", _err))
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
			c.log(ctx, logging.Patch, obj, color.BoldGreen, "OK")
		} else {
			c.log(ctx, logging.Patch, obj, color.BoldYellow, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	return c.inner.Patch(ctx, obj, patch, opts...)
}

func (c *runnerClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return c.inner.IsObjectNamespaced(obj)
}

func (c *runnerClient) log(ctx context.Context, op logging.Operation, obj ctrlclient.Object, color *color.Color, args ...interface{}) {
	logger := logging.FromContext(ctx)
	if logger != nil {
		logger.WithResource(obj).Log(op, color, args...)
	}
}
