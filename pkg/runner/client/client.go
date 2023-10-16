package client

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func New(t *testing.T, inner client.Client, cleanup bool) client.Client {
	t.Helper()
	return &runnerClient{
		t:       t,
		inner:   inner,
		cleanup: cleanup,
	}
}

type runnerClient struct {
	t       *testing.T
	inner   client.Client
	cleanup bool
}

func (c *runnerClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	c.log("create", client.ObjectKey(obj), obj)
	gvk := obj.GetObjectKind().GroupVersionKind()
	err := c.inner.Create(ctx, obj, opts...)
	if err != nil {
		return err
	}
	if c.cleanup {
		c.t.Cleanup(func() {
			obj.GetObjectKind().SetGroupVersionKind(gvk)
			if err := client.BlockingDelete(context.Background(), c, obj); err != nil {
				c.t.Fatal(err)
			}
		})
	}
	return nil
}

func (c *runnerClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	c.log("delete", client.ObjectKey(obj), obj)
	return c.inner.Delete(ctx, obj, opts...)
}

func (c *runnerClient) DeleteAllOf(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteAllOfOption) error {
	c.log("deleteAllOf", client.ObjectKey(obj), obj)
	return c.inner.DeleteAllOf(ctx, obj, opts...)
}

func (c *runnerClient) Get(ctx context.Context, key types.NamespacedName, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	// c.log("get", key, obj)
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
	gvk := obj.GetObjectKind().GroupVersionKind()
	c.t.Logf("%s %s (%s/%s)", op, key, gvk.GroupVersion(), gvk.Kind)
}
