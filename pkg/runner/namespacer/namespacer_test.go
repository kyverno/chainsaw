package namespacer

import (
	"errors"
	"testing"

	mock "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

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
			t.Helper()
			mock := &mock.FakeClient{
				IsNamespaced: func(t *testing.T, obj runtime.Object) (bool, error) {
					t.Helper()
					return tt.namespaced, nil
				},
				ClientErr: tt.clientErr,
			}

			n := New(mock, "test-namespace")
			err := n.Apply(tt.resource)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectNS, tt.resource.GetNamespace())
				assert.Equal(t, tt.expectedCalls, mock.NumCalls)
			}
		})
	}

	t.Run("test GetNamespace", func(t *testing.T) {
		n := New(&mock.FakeClient{}, "test-namespace")
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
