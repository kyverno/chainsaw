package client

import (
	"context"
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		inner client.Client
		want  client.Client
	}{{
		name:  "nil",
		inner: nil,
		want:  &runnerClient{},
	}, {
		name:  "not nil",
		inner: &tclient.FakeClient{},
		want: &runnerClient{
			inner: &tclient.FakeClient{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.inner)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_runnerClient_Get(t *testing.T) {
	type args struct {
		key  types.NamespacedName
		obj  client.Object
		opts []client.GetOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *tlogging.FakeLogger
		inner       func(t *testing.T) *tclient.FakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, _ client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
					return errors.New("test")
				},
			}
		},
		args: args{
			key:  types.NamespacedName{Namespace: "foo", Name: "bar"},
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:    true,
		innerCalls: 1,
	}, {
		name: "no error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, _ client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
					return nil
				},
			}
		},
		args: args{
			key:  types.NamespacedName{Namespace: "foo", Name: "bar"},
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:    false,
		innerCalls: 1,
	}, {
		name: "inner was called",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					assert.Equal(t, types.NamespacedName{Namespace: "foo", Name: "bar"}, key)
					assert.Equal(t, &unstructured.Unstructured{}, obj)
					assert.Nil(t, opts)
					return nil
				},
			}
		},
		args: args{
			key:  types.NamespacedName{Namespace: "foo", Name: "bar"},
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:    false,
		innerCalls: 1,
	}, {
		name: "logger was not called",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					assert.Equal(t, types.NamespacedName{Namespace: "foo", Name: "bar"}, key)
					assert.Equal(t, &unstructured.Unstructured{}, obj)
					assert.Nil(t, opts)
					return nil
				},
			}
		},
		args: args{
			key:  types.NamespacedName{Namespace: "foo", Name: "bar"},
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:    false,
		innerCalls: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			ctx := logging.IntoContext(context.TODO(), mockLogger)
			err := c.Get(ctx, tt.args.key, tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
			assert.Equal(t, tt.loggerCalls, mockLogger.NumCalls())
		})
	}
}

func Test_runnerClient_Create(t *testing.T) {
	type args struct {
		obj  client.Object
		opts []client.CreateOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *tlogging.FakeLogger
		inner       func(t *testing.T) *tclient.FakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
					return errors.New("test")
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     true,
		loggerCalls: 2,
		innerCalls:  1,
	}, {
		name: "no error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
					return nil
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     false,
		loggerCalls: 2,
		innerCalls:  1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			ctx := logging.IntoContext(context.TODO(), mockLogger)
			err := c.Create(ctx, tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
			assert.Equal(t, tt.loggerCalls, mockLogger.NumCalls())
		})
	}
}

func Test_runnerClient_Update(t *testing.T) {
	type args struct {
		obj  client.Object
		opts []client.UpdateOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *tlogging.FakeLogger
		inner       func(t *testing.T) *tclient.FakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				UpdateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.UpdateOption) error {
					return errors.New("test")
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     true,
		loggerCalls: 2,
		innerCalls:  1,
	}, {
		name: "no error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				UpdateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.UpdateOption) error {
					return nil
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     false,
		loggerCalls: 2,
		innerCalls:  1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			ctx := logging.IntoContext(context.TODO(), mockLogger)
			err := c.Update(ctx, tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
			assert.Equal(t, tt.loggerCalls, mockLogger.NumCalls())
		})
	}
}

func Test_runnerClient_Delete(t *testing.T) {
	type args struct {
		obj  client.Object
		opts []client.DeleteOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *tlogging.FakeLogger
		inner       func(t *testing.T) *tclient.FakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				DeleteFn: func(_ context.Context, _ int, _ client.Object, _ ...client.DeleteOption) error {
					return errors.New("test")
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     true,
		loggerCalls: 2,
		innerCalls:  1,
	}, {
		name: "no error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				DeleteFn: func(_ context.Context, _ int, _ client.Object, _ ...client.DeleteOption) error {
					return nil
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     false,
		loggerCalls: 2,
		innerCalls:  1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			ctx := logging.IntoContext(context.TODO(), mockLogger)
			err := c.Delete(ctx, tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
			assert.Equal(t, tt.loggerCalls, mockLogger.NumCalls())
		})
	}
}

func Test_runnerClient_List(t *testing.T) {
	type args struct {
		obj  client.ObjectList
		opts []client.ListOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *tlogging.FakeLogger
		inner       func(t *testing.T) *tclient.FakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				ListFn: func(_ context.Context, _ int, _ client.ObjectList, _ ...client.ListOption) error {
					return errors.New("test")
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     true,
		loggerCalls: 0,
		innerCalls:  1,
	}, {
		name: "no error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				ListFn: func(_ context.Context, _ int, _ client.ObjectList, _ ...client.ListOption) error {
					return nil
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     false,
		loggerCalls: 0,
		innerCalls:  1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			ctx := logging.IntoContext(context.TODO(), mockLogger)
			err := c.List(ctx, tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
			assert.Equal(t, tt.loggerCalls, mockLogger.NumCalls())
		})
	}
}

func Test_runnerClient_Patch(t *testing.T) {
	type args struct {
		obj   client.Object
		patch client.Patch
		opts  []client.PatchOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *tlogging.FakeLogger
		inner       func(t *testing.T) *tclient.FakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
					return errors.New("test")
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     true,
		loggerCalls: 2,
		innerCalls:  1,
	}, {
		name: "no error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
					return nil
				},
			}
		},
		args: args{
			obj:  &unstructured.Unstructured{},
			opts: nil,
		},
		wantErr:     false,
		loggerCalls: 2,
		innerCalls:  1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			ctx := logging.IntoContext(context.TODO(), mockLogger)
			err := c.Patch(ctx, tt.args.obj, tt.args.patch, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
			assert.Equal(t, tt.loggerCalls, mockLogger.NumCalls())
		})
	}
}

func Test_runnerClient_IsObjectNamespaced(t *testing.T) {
	type args struct {
		obj runtime.Object
	}
	tests := []struct {
		name       string
		logger     func(t *testing.T) *tlogging.FakeLogger
		inner      func(t *testing.T) *tclient.FakeClient
		args       args
		want       bool
		wantErr    bool
		innerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				IsObjectNamespacedFn: func(_ int, _ runtime.Object) (bool, error) {
					return false, errors.New("test")
				},
			}
		},
		args: args{
			obj: &unstructured.Unstructured{},
		},
		wantErr:    true,
		innerCalls: 1,
	}, {
		name: "no error - false",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				IsObjectNamespacedFn: func(_ int, _ runtime.Object) (bool, error) {
					return false, nil
				},
			}
		},
		args: args{
			obj: &unstructured.Unstructured{},
		},
		want:       false,
		wantErr:    false,
		innerCalls: 1,
	}, {
		name: "no error - true",
		logger: func(t *testing.T) *tlogging.FakeLogger {
			t.Helper()
			return &tlogging.FakeLogger{}
		},
		inner: func(t *testing.T) *tclient.FakeClient {
			t.Helper()
			return &tclient.FakeClient{
				IsObjectNamespacedFn: func(_ int, _ runtime.Object) (bool, error) {
					return true, nil
				},
			}
		},
		args: args{
			obj: &unstructured.Unstructured{},
		},
		want:       true,
		wantErr:    false,
		innerCalls: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.inner(t)
			c := &runnerClient{
				inner: mockClient,
			}
			got, err := c.IsObjectNamespaced(tt.args.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls())
		})
	}
}

func Test_runnerClient_RESTMapper(t *testing.T) {
	inner := &tclient.FakeClient{
		RESTMapperFn: func(_ int) meta.RESTMapper {
			return nil
		},
	}
	c := &runnerClient{
		inner: inner,
	}
	got := c.RESTMapper()
	assert.Equal(t, 1, inner.NumCalls())
	assert.Nil(t, got)
}
