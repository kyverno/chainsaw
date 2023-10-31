package namespacer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type mockClient struct {
	t            *testing.T
	clientErr    error
	get          func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error
	create       func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error
	delete       func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error
	list         func(ctx context.Context, t *testing.T, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error
	patch        func(ctx context.Context, t *testing.T, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error
	isNamespaced func(t *testing.T, obj runtime.Object) (bool, error)
	numCalls     int
}

func (f *mockClient) Get(ctx context.Context, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
	defer func() { f.numCalls++ }()
	return f.get(ctx, f.t, key, obj, opts...)
}

func (f *mockClient) List(ctx context.Context, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
	defer func() { f.numCalls++ }()
	return f.list(ctx, f.t, list, opts...)
}

func (f *mockClient) Create(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
	defer func() { f.numCalls++ }()
	return f.create(ctx, f.t, obj, opts...)
}

func (f *mockClient) Delete(ctx context.Context, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
	defer func() { f.numCalls++ }()
	return f.delete(ctx, f.t, obj, opts...)
}

func (f *mockClient) Patch(ctx context.Context, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
	defer func() { f.numCalls++ }()
	return f.patch(ctx, f.t, obj, patch, opts...)
}

func (f *mockClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	defer func() { f.numCalls++ }()
	if f.clientErr != nil {
		return false, f.clientErr
	}
	return f.isNamespaced(f.t, obj)
}

func TestNamespacer(t *testing.T) {
	tests := []struct {
		name          string
		resource      ctrlclient.Object
		namespaced    bool
		expectErr     bool
		expectNS      string
		clientErr     error
		expectedCalls int
	}{
		{
			name:      "nil resource",
			resource:  nil,
			expectErr: true,
		},
		{
			name:          "namespaced with no NS set",
			resource:      &testResource{},
			namespaced:    true,
			expectNS:      "test-namespace",
			expectedCalls: 1,
		},
		{
			name:          "non-namespaced resource",
			resource:      &testResource{},
			expectNS:      "",
			expectedCalls: 1,
		},
		{
			name:     "resource with namespace set",
			resource: &testResource{namespace: "already-set"},
			expectNS: "already-set",
		},
		{
			name:          "error scenario",
			resource:      &testResource{},
			clientErr:     errors.New("mock error"),
			expectErr:     true,
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockClient{
				isNamespaced: func(t *testing.T, obj runtime.Object) (bool, error) {
					return tt.namespaced, nil
				},
				clientErr: tt.clientErr,
			}

			n := New(mock, "test-namespace")
			err := n.Apply(tt.resource)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectNS, tt.resource.GetNamespace())
				assert.Equal(t, tt.expectedCalls, mock.numCalls)
			}
		})
	}

	t.Run("test GetNamespace", func(t *testing.T) {
		n := New(&mockClient{}, "test-namespace")
		ns := n.GetNamespace()
		assert.Equal(t, "test-namespace", ns)
	})
}

// Mock the ctrlclient.Object for testing purposes.
type testResource struct {
	ctrlclient.Object
	namespace string
}

func (t *testResource) GetNamespace() string {
	return t.namespace
}

func (t *testResource) SetNamespace(ns string) {
	t.namespace = ns
}
