package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFakeClient(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		called := false
		c := &FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		err := c.Get(context.Background(), client.ObjectKey{}, nil)
		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("List", func(t *testing.T) {
		called := false
		c := &FakeClient{
			ListFn: func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		err := c.List(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("Create", func(t *testing.T) {
		called := false
		c := &FakeClient{
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		err := c.Create(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("Update", func(t *testing.T) {
		called := false
		c := &FakeClient{
			UpdateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.UpdateOption) error {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		err := c.Update(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("Delete", func(t *testing.T) {
		called := false
		c := &FakeClient{
			DeleteFn: func(ctx context.Context, call int, obj client.Object, opts ...client.DeleteOption) error {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		err := c.Delete(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("Patch", func(t *testing.T) {
		called := false
		c := &FakeClient{
			PatchFn: func(ctx context.Context, call int, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		err := c.Patch(context.Background(), nil, nil)
		assert.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("IsObjectNamespaced", func(t *testing.T) {
		called := false
		c := &FakeClient{
			IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
				called = true
				assert.Equal(t, 0, call)
				return true, nil
			},
		}
		namespaced, err := c.IsObjectNamespaced(nil)
		assert.NoError(t, err)
		assert.True(t, namespaced)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("RESTMapper", func(t *testing.T) {
		called := false
		c := &FakeClient{
			RESTMapperFn: func(call int) meta.RESTMapper {
				called = true
				assert.Equal(t, 0, call)
				return nil
			},
		}
		mapper := c.RESTMapper()
		assert.Nil(t, mapper)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("SubResource - with Fn", func(t *testing.T) {
		called := false
		c := &FakeClient{
			SubResourceFn: func(subResource string) client.SubResourceClient {
				called = true
				assert.Equal(t, "status", subResource)
				return nil
			},
		}
		sw := c.SubResource("status")
		assert.Nil(t, sw)
		assert.True(t, called)
		assert.Equal(t, 1, c.NumCalls())
	})

	t.Run("SubResource - without Fn", func(t *testing.T) {
		c := &FakeClient{}
		sw := c.SubResource("status")
		assert.NotNil(t, sw)
		assert.Equal(t, 1, c.NumCalls())
	})
}

func TestFakeSubResourceWriter(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		called := false
		f := &FakeSubResourceWriter{
			GetFn: func(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceGetOption) error {
				called = true
				return nil
			},
		}
		err := f.Get(context.Background(), nil, nil)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Create", func(t *testing.T) {
		called := false
		f := &FakeSubResourceWriter{
			CreateFn: func(ctx context.Context, obj client.Object, subResource client.Object, opts ...client.SubResourceCreateOption) error {
				called = true
				return nil
			},
		}
		err := f.Create(context.Background(), nil, nil)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Update", func(t *testing.T) {
		called := false
		f := &FakeSubResourceWriter{
			UpdateFn: func(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
				called = true
				return nil
			},
		}
		err := f.Update(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Patch", func(t *testing.T) {
		called := false
		f := &FakeSubResourceWriter{
			PatchFn: func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
				called = true
				return nil
			},
		}
		err := f.Patch(context.Background(), nil, nil)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Apply", func(t *testing.T) {
		called := false
		f := &FakeSubResourceWriter{
			ApplyFn: func(ctx context.Context, obj runtime.ApplyConfiguration, opts ...client.SubResourceApplyOption) error {
				called = true
				return nil
			},
		}
		err := f.Apply(context.Background(), nil)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("NewFakeSubResourceWriter", func(t *testing.T) {
		f := NewFakeSubResourceWriter()
		assert.NotNil(t, f)
		assert.NoError(t, f.Update(context.Background(), nil))
		assert.NoError(t, f.Patch(context.Background(), nil, nil))
		assert.NoError(t, f.Create(context.Background(), nil, nil))

		// Get and Apply are nil by default in NewFakeSubResourceWriter
		assert.Panics(t, func() { _ = f.Get(context.Background(), nil, nil) })
		assert.Panics(t, func() { _ = f.Apply(context.Background(), nil) })
	})
}
