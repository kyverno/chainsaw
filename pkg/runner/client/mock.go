package client

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type FakeClient struct {
	t            *testing.T
	ClientErr    error
	get          func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error
	create       func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error
	delete       func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error
	list         func(ctx context.Context, t *testing.T, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error
	patch        func(ctx context.Context, t *testing.T, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error
	IsNamespaced func(t *testing.T, obj runtime.Object) (bool, error)
	NumCalls     int
}

func (f *FakeClient) Get(ctx context.Context, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	defer func() { f.NumCalls++ }()
	return f.get(ctx, f.t, key, obj, opts...)
}

func (f *FakeClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
	defer func() { f.NumCalls++ }()
	return f.list(ctx, f.t, list, opts...)
}

func (f *FakeClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	defer func() { f.NumCalls++ }()
	return f.create(ctx, f.t, obj, opts...)
}

func (f *FakeClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	defer func() { f.NumCalls++ }()
	return f.delete(ctx, f.t, obj, opts...)
}

func (f *FakeClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
	defer func() { f.NumCalls++ }()
	return f.patch(ctx, f.t, obj, patch, opts...)
}

func (f *FakeClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	defer func() { f.NumCalls++ }()
	if f.ClientErr != nil {
		return false, f.ClientErr
	}
	return f.IsNamespaced(f.t, obj)
}
