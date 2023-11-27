package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
)

func TestToUnstructured(t *testing.T) {
	assert.Panics(t, func() {
		ToUnstructured(nil)
	})
	assert.Equal(t, ToUnstructured(ptr.To(Namespace("foo"))), unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				"creationTimestamp": nil,
				"name":              "foo",
			},
			"spec":   map[string]interface{}{},
			"status": map[string]interface{}{},
		},
	})
}
