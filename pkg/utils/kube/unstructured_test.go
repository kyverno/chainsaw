package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestToUnstructured(t *testing.T) {
	var nilPtr *int
	assert.Panics(t, func() {
		ToUnstructured(nilPtr)
	})
	assert.Equal(t, ToUnstructured(new(Namespace("foo"))), unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]any{
				"name": "foo",
			},
			"spec":   map[string]any{},
			"status": map[string]any{},
		},
	})
}
