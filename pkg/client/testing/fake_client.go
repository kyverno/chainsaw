package testing

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: not thread safe
type FakeClient struct {
	GetFn          func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error
	CreateFn       func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error
	DeleteFn       func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error
	ListFn         func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error
	PatchFn        func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error
	IsNamespacedFn func(call int, obj runtime.Object) (bool, error)
	numCalls       int
	StatusWriterFn func(call int) ctrlclient.StatusWriter
}

func (c *FakeClient) Get(ctx context.Context, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	defer func() { c.numCalls++ }()
	return c.GetFn(ctx, c.numCalls, key, obj, opts...)
}

func (c *FakeClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
	defer func() { c.numCalls++ }()
	return c.ListFn(ctx, c.numCalls, list, opts...)
}

func (c *FakeClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	defer func() { c.numCalls++ }()
	return c.CreateFn(ctx, c.numCalls, obj, opts...)
}

func (c *FakeClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	defer func() { c.numCalls++ }()
	return c.DeleteFn(ctx, c.numCalls, obj, opts...)
}

func (c *FakeClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
	defer func() { c.numCalls++ }()
	return c.PatchFn(ctx, c.numCalls, obj, patch, opts...)
}

func (c *FakeClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	defer func() { c.numCalls++ }()
	return c.IsNamespacedFn(c.numCalls, obj)
}

func (c *FakeClient) NumCalls() int {
	return c.numCalls
}

type FakeStatusWriter struct {
	UpdateFn func(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.SubResourceUpdateOption) error
	PatchFn  func(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.SubResourcePatchOption) error
	CreateFn func(ctx context.Context, obj ctrlclient.Object, subResource ctrlclient.Object, opts ...ctrlclient.SubResourceCreateOption) error
}

func (f *FakeStatusWriter) Update(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.SubResourceUpdateOption) error {
	return f.UpdateFn(ctx, obj, opts...)
}

func (f *FakeStatusWriter) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.SubResourcePatchOption) error {
	return f.PatchFn(ctx, obj, patch, opts...)
}

func (f *FakeStatusWriter) Create(ctx context.Context, obj ctrlclient.Object, subResource ctrlclient.Object, opts ...ctrlclient.SubResourceCreateOption) error {
	return f.CreateFn(ctx, obj, subResource, opts...)
}

func NewFakeStatusWriter() *FakeStatusWriter {
	return &FakeStatusWriter{
		UpdateFn: func(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.SubResourceUpdateOption) error {
			return nil
		},
		PatchFn: func(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.SubResourcePatchOption) error {
			return nil
		},
		CreateFn: func(ctx context.Context, obj ctrlclient.Object, subResource ctrlclient.Object, opts ...ctrlclient.SubResourceCreateOption) error {
			return nil
		},
	}
}

func (c *FakeClient) Status() ctrlclient.StatusWriter {
	if c.StatusWriterFn != nil {
		return c.StatusWriterFn(c.numCalls)
	}
	return NewFakeStatusWriter()
}
