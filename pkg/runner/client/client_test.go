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
		inner:  &FakeClient{},
		want: &runnerClient{
			inner: &FakeClient{},
		},
	}, {
		name:   "logger and client",
		logger: &fakeLogger{},
		inner:  &FakeClient{},
		want: &runnerClient{
			logger: &fakeLogger{},
			inner:  &FakeClient{},
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
		inner       func(t *testing.T) *FakeClient
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				GetFake: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				GetFake: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				GetFake: func(_ context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				GetFake: func(_ context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
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
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls)
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
		inner       func(t *testing.T) *FakeClient
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
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
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls)
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
		inner       func(t *testing.T) *FakeClient
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
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
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls)
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
		inner       func(t *testing.T) *FakeClient
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				ListFake: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectList, _ ...ctrlclient.ListOption) error {
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				ListFake: func(_ context.Context, t *testing.T, _ ctrlclient.ObjectList, _ ...ctrlclient.ListOption) error {
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
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls)
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
		inner       func(t *testing.T) *FakeClient
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
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
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls)
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
		inner       func(t *testing.T) *FakeClient
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				IsNamespaced: func(t *testing.T, _ runtime.Object) (bool, error) {
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				IsNamespaced: func(t *testing.T, _ runtime.Object) (bool, error) {
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
		inner: func(t *testing.T) *FakeClient {
			t.Helper()
			return &FakeClient{
				T: t,
				IsNamespaced: func(t *testing.T, _ runtime.Object) (bool, error) {
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
			assert.Equal(t, tt.innerCalls, mockClient.NumCalls)
			assert.Equal(t, tt.loggerCalls, mockLogger.numCalls)
		})
	}
}
