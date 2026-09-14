package testing

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

type FakeClient struct {
	mu                   sync.RWMutex
	GetFn                func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
	CreateFn             func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error
	UpdateFn             func(ctx context.Context, call int, obj client.Object, opts ...client.UpdateOption) error
	DeleteFn             func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error
	ListFn               func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error
	PatchFn              func(ctx context.Context, call int, obj client.Object, patch client.Patch, opts ...client.PatchOption) error
	SubResourceFn        func(subResource string) client.SubResourceClient
	IsObjectNamespacedFn func(call int, obj runtime.Object) (bool, error)
	RESTMapperFn         func(call int) meta.RESTMapper
	numCalls             int
}

func (c *FakeClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.GetFn(ctx, call, key, obj, opts...)
}

func (c *FakeClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.ListFn(ctx, call, list, opts...)
}

func (c *FakeClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.CreateFn(ctx, call, obj, opts...)
}

func (c *FakeClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.UpdateFn(ctx, call, obj, opts...)
}

func (c *FakeClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.DeleteFn(ctx, call, obj, opts...)
}

func (c *FakeClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.PatchFn(ctx, call, obj, patch, opts...)
}

func (c *FakeClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.IsObjectNamespacedFn(call, obj)
}

func (c *FakeClient) RESTMapper() meta.RESTMapper {
	c.mu.Lock()
	call := c.numCalls
	c.numCalls++
	c.mu.Unlock()
	return c.RESTMapperFn(call)
}

func (c *FakeClient) SubResource(subResource string) client.SubResourceClient {
	c.mu.Lock()
	c.numCalls++
	c.mu.Unlock()
	if c.SubResourceFn != nil {
		return c.SubResourceFn(subResource)
	}
	return NewFakeSubResourceWriter()
}

func (c *FakeClient) NumCalls() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.numCalls
}

type FakeSubResourceWriter struct {
	GetFn    func(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceGetOption) error
	CreateFn func(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error
	UpdateFn func(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error
	PatchFn  func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error
	ApplyFn  func(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.SubResourceApplyOption) error
}

func (f *FakeSubResourceWriter) Get(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceGetOption) error {
	return f.GetFn(ctx, obj, subResource, opts...)
}

func (f *FakeSubResourceWriter) Create(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	return f.CreateFn(ctx, obj, subResource, opts...)
}

func (f *FakeSubResourceWriter) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	return f.UpdateFn(ctx, obj, opts...)
}

func (f *FakeSubResourceWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	return f.PatchFn(ctx, obj, patch, opts...)
}

func (f *FakeSubResourceWriter) Apply(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.SubResourceApplyOption) error {
	return f.ApplyFn(ctx, obj, opts...)
}

func NewFakeSubResourceWriter() *FakeSubResourceWriter {
	return &FakeSubResourceWriter{
		UpdateFn: func(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
			return nil
		},
		PatchFn: func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
			return nil
		},
		CreateFn: func(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error {
			return nil
		},
	}
}
