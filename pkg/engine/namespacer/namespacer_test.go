package namespacer

import (
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestNamespacer(t *testing.T) {
	var nsNotSet unstructured.Unstructured
	var nsSet unstructured.Unstructured
	nsSet.SetNamespace("already-set")
	tests := []struct {
		name          string
		resource      client.Object
		namespaced    bool
		expectErr     bool
		expectNS      string
		clientErr     error
		expectedCalls int
	}{{
		name:      "nil resource",
		resource:  nil,
		expectErr: true,
	}, {
		name:          "namespaced with no NS set",
		resource:      nsNotSet.DeepCopy(),
		namespaced:    true,
		expectNS:      "test-namespace",
		expectedCalls: 1,
	}, {
		name:          "non-namespaced resource",
		namespaced:    false,
		resource:      nsNotSet.DeepCopy(),
		expectNS:      "",
		expectedCalls: 1,
	}, {
		name:     "resource with namespace set",
		resource: nsSet.DeepCopy(),
		expectNS: "already-set",
	}, {
		name:          "error scenario",
		resource:      nsNotSet.DeepCopy(),
		clientErr:     errors.New("mock error"),
		expectErr:     true,
		expectedCalls: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &tclient.FakeClient{
				IsObjectNamespacedFn: func(_ int, obj runtime.Object) (bool, error) {
					t.Helper()
					if tt.clientErr != nil {
						return false, tt.clientErr
					}
					return tt.namespaced, nil
				},
			}
			n := New("test-namespace")
			err := n.Apply(mock, tt.resource)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectNS, tt.resource.GetNamespace())
				assert.Equal(t, tt.expectedCalls, mock.NumCalls())
			}
		})
	}
	t.Run("test GetNamespace", func(t *testing.T) {
		n := New("test-namespace")
		ns := n.GetNamespace()
		assert.Equal(t, "test-namespace", ns)
	})
}
