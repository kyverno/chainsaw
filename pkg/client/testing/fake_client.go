package testing

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: not thread safe
type FakeClient struct {
	GetFn                func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
	CreateFn             func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error
	UpdateFn             func(ctx context.Context, call int, obj client.Object, opts ...client.UpdateOption) error
	DeleteFn             func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error
	ListFn               func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error
	PatchFn              func(ctx context.Context, call int, obj client.Object, patch client.Patch, opts ...client.PatchOption) error
	IsObjectNamespacedFn func(call int, obj runtime.Object) (bool, error)
	RESTMapperFn         func(call int) meta.RESTMapper
	numCalls             int
}

func (c *FakeClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	defer func() { c.numCalls++ }()
	return c.GetFn(ctx, c.numCalls, key, obj, opts...)
}

func (c *FakeClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	defer func() { c.numCalls++ }()
	return c.ListFn(ctx, c.numCalls, list, opts...)
}

func (c *FakeClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	defer func() { c.numCalls++ }()
	return c.CreateFn(ctx, c.numCalls, obj, opts...)
}

func (c *FakeClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	defer func() { c.numCalls++ }()
	return c.UpdateFn(ctx, c.numCalls, obj, opts...)
}

func (c *FakeClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	defer func() { c.numCalls++ }()
	return c.DeleteFn(ctx, c.numCalls, obj, opts...)
}

func (c *FakeClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	defer func() { c.numCalls++ }()
	return c.PatchFn(ctx, c.numCalls, obj, patch, opts...)
}

func (c *FakeClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	defer func() { c.numCalls++ }()
	return c.IsObjectNamespacedFn(c.numCalls, obj)
}

func (c *FakeClient) RESTMapper() meta.RESTMapper {
	defer func() { c.numCalls++ }()
	return c.RESTMapperFn(c.numCalls)
}

func (c *FakeClient) NumCalls() int {
	return c.numCalls
}
