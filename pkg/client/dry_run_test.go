package client

import (
	"context"
	"errors"
	"testing"

	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestDryRun(t *testing.T) {
	tests := []struct {
		name  string
		inner Client
		want  Client
	}{{
		name:  "nil",
		inner: nil,
		want: &dryRunClient{
			inner: nil,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DryRun(tt.inner)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_dryRunClient_Create(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.Object
		opts    []client.CreateOption
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return false, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			err := c.Create(context.TODO(), tt.obj, tt.opts...)
			assert.Equal(t, 1, inner.NumCalls())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_Update(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.Object
		opts    []client.UpdateOption
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return false, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			err := c.Update(context.TODO(), tt.obj, tt.opts...)
			assert.Equal(t, 1, inner.NumCalls())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_Delete(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.Object
		opts    []client.DeleteOption
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return false, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			err := c.Delete(context.TODO(), tt.obj, tt.opts...)
			assert.Equal(t, 1, inner.NumCalls())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_Get(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.Object
		opts    []client.GetOption
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return false, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			err := c.Get(context.TODO(), types.NamespacedName{}, tt.obj, tt.opts...)
			assert.Equal(t, 1, inner.NumCalls())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_List(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.ObjectList
		opts    []client.ListOption
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return false, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			err := c.List(context.TODO(), tt.obj, tt.opts...)
			assert.Equal(t, 1, inner.NumCalls())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_Patch(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.Object
		opts    []client.PatchOption
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return false, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			err := c.Patch(context.TODO(), tt.obj, nil, tt.opts...)
			assert.Equal(t, 1, inner.NumCalls())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_IsObjectNamespaced(t *testing.T) {
	tests := []struct {
		name    string
		inner   Client
		obj     client.Object
		opts    []client.PatchOption
		want    bool
		wantErr bool
	}{{
		name:    "no error",
		obj:     nil,
		opts:    nil,
		want:    false,
		wantErr: false,
	}, {
		name:    "error",
		obj:     nil,
		opts:    nil,
		want:    false,
		wantErr: true,
	}, {
		name:    "true",
		obj:     nil,
		opts:    nil,
		want:    true,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr := func() error {
				if tt.wantErr {
					return errors.New("dummy error")
				}
				return nil
			}
			inner := &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				UpdateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.UpdateOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					assert.NotContains(t, opts, client.DryRunAll)
					return wantErr()
				},
				PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
					assert.Contains(t, opts, client.DryRunAll)
					return wantErr()
				},
				IsObjectNamespacedFn: func(call int, obj runtime.Object) (bool, error) {
					return tt.want, wantErr()
				},
			}
			c := &dryRunClient{
				inner: inner,
			}
			got, err := c.IsObjectNamespaced(tt.obj)
			assert.Equal(t, 1, inner.NumCalls())
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_dryRunClient_RESTMapper(t *testing.T) {
	inner := &tclient.FakeClient{
		RESTMapperFn: func(call int) meta.RESTMapper {
			return nil
		},
	}
	c := &dryRunClient{
		inner: inner,
	}
	got := c.RESTMapper()
	assert.Equal(t, 1, inner.NumCalls())
	assert.Nil(t, got)
}
