package client

import (
	"context"
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type fakeLogger struct {
	log      func(int, string, *color.Color, ...interface{})
	numCalls int
}

func (f *fakeLogger) Log(s string, c *color.Color, a ...interface{}) {
	defer func() { f.numCalls++ }()
	if f.log != nil {
		f.log(f.numCalls, s, c, a...)
	}
}

func (f *fakeLogger) WithResource(ctrlclient.Object) logging.Logger {
	defer func() { f.numCalls++ }()
	return f
}

type fakeClient struct {
	t            *testing.T
	get          func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error
	create       func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error
	delete       func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error
	list         func(ctx context.Context, t *testing.T, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error
	patch        func(ctx context.Context, t *testing.T, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error
	isNamespaced func(t *testing.T, obj runtime.Object) (bool, error)
	numCalls     int
}

func (f *fakeClient) Get(ctx context.Context, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	defer func() { f.numCalls++ }()
	return f.get(ctx, f.t, key, obj, opts...)
}

func (f *fakeClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
	defer func() { f.numCalls++ }()
	return f.list(ctx, f.t, list, opts...)
}

func (f *fakeClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	defer func() { f.numCalls++ }()
	return f.create(ctx, f.t, obj, opts...)
}

func (f *fakeClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	defer func() { f.numCalls++ }()
	return f.delete(ctx, f.t, obj, opts...)
}

func (f *fakeClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
	defer func() { f.numCalls++ }()
	return f.patch(ctx, f.t, obj, patch, opts...)
}

func (f *fakeClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	defer func() { f.numCalls++ }()
	return f.isNamespaced(f.t, obj)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		logger logging.Logger
		inner  client.Client
		want   client.Client
	}{{
		name:   "nil",
		logger: nil,
		inner:  nil,
		want:   &runnerClient{},
	}, {
		name:   "only logger",
		logger: &fakeLogger{},
		inner:  nil,
		want: &runnerClient{
			logger: &fakeLogger{},
		},
	}, {
		name:   "only client",
		logger: nil,
		inner:  &fakeClient{},
		want: &runnerClient{
			inner: &fakeClient{},
		},
	}, {
		name:   "logger and client",
		logger: &fakeLogger{},
		inner:  &fakeClient{},
		want: &runnerClient{
			logger: &fakeLogger{},
			inner:  &fakeClient{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.logger, tt.inner)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_runnerClient_Get(t *testing.T) {
	type args struct {
		key  types.NamespacedName
		obj  ctrlclient.Object
		opts []ctrlclient.GetOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *fakeLogger
		inner       func(t *testing.T) *fakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				get: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				get: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				get: func(_ context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				get: func(_ context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					t.Helper()
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
				logger: mockLogger,
				inner:  mockClient,
			}
			err := c.Get(context.TODO(), tt.args.key, tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.numCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}

func Test_runnerClient_Create(t *testing.T) {
	type args struct {
		obj  ctrlclient.Object
		opts []ctrlclient.CreateOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *fakeLogger
		inner       func(t *testing.T) *fakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				create: func(_ context.Context, t *testing.T, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				create: func(_ context.Context, t *testing.T, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					t.Helper()
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
				logger: mockLogger,
				inner:  mockClient,
			}
			err := c.Create(context.TODO(), tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.numCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}

func Test_runnerClient_Delete(t *testing.T) {
	type args struct {
		obj  ctrlclient.Object
		opts []ctrlclient.DeleteOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *fakeLogger
		inner       func(t *testing.T) *fakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				delete: func(_ context.Context, t *testing.T, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				delete: func(_ context.Context, t *testing.T, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					t.Helper()
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
				logger: mockLogger,
				inner:  mockClient,
			}
			err := c.Delete(context.TODO(), tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.numCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}

func Test_runnerClient_List(t *testing.T) {
	type args struct {
		obj  ctrlclient.ObjectList
		opts []ctrlclient.ListOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *fakeLogger
		inner       func(t *testing.T) *fakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				list: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectList, _ ...ctrlclient.ListOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				list: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectList, _ ...ctrlclient.ListOption) error {
					t.Helper()
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
				logger: mockLogger,
				inner:  mockClient,
			}
			err := c.List(context.TODO(), tt.args.obj, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.numCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}

func Test_runnerClient_Patch(t *testing.T) {
	type args struct {
		obj   ctrlclient.Object
		patch ctrlclient.Patch
		opts  []ctrlclient.PatchOption
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *fakeLogger
		inner       func(t *testing.T) *fakeClient
		args        args
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				patch: func(_ context.Context, t *testing.T, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
					t.Helper()
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
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				patch: func(_ context.Context, t *testing.T, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
					t.Helper()
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
				logger: mockLogger,
				inner:  mockClient,
			}
			err := c.Patch(context.TODO(), tt.args.obj, tt.args.patch, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.innerCalls, mockClient.numCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}

func Test_runnerClient_IsObjectNamespaced(t *testing.T) {
	type args struct {
		obj runtime.Object
	}
	tests := []struct {
		name        string
		logger      func(t *testing.T) *fakeLogger
		inner       func(t *testing.T) *fakeClient
		args        args
		want        bool
		wantErr     bool
		innerCalls  int
		loggerCalls int
	}{{
		name: "with error",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				isNamespaced: func(t *testing.T, _ runtime.Object) (bool, error) {
					t.Helper()
					return false, errors.New("test")
				},
			}
		},
		args: args{
			obj: &unstructured.Unstructured{},
		},
		wantErr:     true,
		loggerCalls: 0,
		innerCalls:  1,
	}, {
		name: "no error - false",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				isNamespaced: func(t *testing.T, _ runtime.Object) (bool, error) {
					t.Helper()
					return false, nil
				},
			}
		},
		args: args{
			obj: &unstructured.Unstructured{},
		},
		want:        false,
		wantErr:     false,
		loggerCalls: 0,
		innerCalls:  1,
	}, {
		name: "no error - true",
		logger: func(t *testing.T) *fakeLogger {
			t.Helper()
			return &fakeLogger{}
		},
		inner: func(t *testing.T) *fakeClient {
			t.Helper()
			return &fakeClient{
				t: t,
				isNamespaced: func(t *testing.T, _ runtime.Object) (bool, error) {
					t.Helper()
					return true, nil
				},
			}
		},
		args: args{
			obj: &unstructured.Unstructured{},
		},
		want:        true,
		wantErr:     false,
		loggerCalls: 0,
		innerCalls:  1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.logger(t)
			mockClient := tt.inner(t)
			c := &runnerClient{
				logger: mockLogger,
				inner:  mockClient,
			}
			got, err := c.IsObjectNamespaced(tt.args.obj)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.Equal(t, tt.innerCalls, mockClient.numCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}
