package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	truntime "k8s.io/apimachinery/pkg/runtime/testing"
)

func TestPatchObject(t *testing.T) {
	tests := []struct {
		name     string
		actual   runtime.Object
		expected runtime.Object
		want     runtime.Object
		wantErr  bool
	}{{
		name:     "actual nil",
		actual:   nil,
		expected: &unstructured.Unstructured{},
		wantErr:  true,
	}, {
		name:     "expected nil",
		actual:   &unstructured.Unstructured{},
		expected: nil,
		wantErr:  true,
	}, {
		name: "ok",
		actual: &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":            "test-pod",
					"resourceVersion": "12345",
				},
			},
		},
		expected: &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name": "test-pod",
				},
				"foo": "bar",
			},
		},
		want: &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":            "test-pod",
					"resourceVersion": "12345",
				},
				"foo": "bar",
			},
		},
	}, {
		name:     "actual not meta",
		actual:   &truntime.InternalSimple{},
		expected: &unstructured.Unstructured{},
		wantErr:  true,
	}, {
		name:     "expected not meta",
		actual:   &unstructured.Unstructured{},
		expected: &truntime.InternalSimple{},
		wantErr:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PatchObject(tt.actual, tt.expected)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
