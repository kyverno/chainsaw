package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestPrettyDiff(t *testing.T) {
	var pod unstructured.Unstructured
	pod.SetAPIVersion("v1")
	pod.SetKind("Pod")
	pod.SetName("foo")
	pod.SetNamespace("bar")
	pod.SetAnnotations(map[string]string{"bar": "baz"})
	assert.NoError(t, unstructured.SetNestedSlice(pod.UnstructuredContent(), []any{"foo"}, "spec", "data"))
	assert.NoError(t, unstructured.SetNestedMap(pod.UnstructuredContent(), map[string]any{"foo": "foos"}, "spec", "something"))
	tests := []struct {
		name     string
		expected unstructured.Unstructured
		actual   unstructured.Unstructured
		want     string
		wantErr  bool
	}{{
		name: "empty",
	}, {
		name: "error",
		expected: unstructured.Unstructured{
			Object: map[string]interface{}{
				"": func() {},
			},
		},
		wantErr: true,
	}, {
		name: "error",
		expected: unstructured.Unstructured{
			Object: map[string]interface{}{
				"": "",
			},
		},
		actual: unstructured.Unstructured{
			Object: map[string]interface{}{
				"": func() {},
			},
		},
		wantErr: true,
	}, {
		name:     "same",
		expected: pod,
		actual:   pod,
	}, {
		name:     "same",
		expected: pod,
		actual: func() unstructured.Unstructured {
			pod := pod.DeepCopy()
			pod.SetResourceVersion("123")
			assert.NoError(t, unstructured.SetNestedMap(pod.UnstructuredContent(), map[string]any{"bar": "bars"}, "spec", "something", "else"))
			assert.NoError(t, unstructured.SetNestedSlice(pod.UnstructuredContent(), []any{"bar"}, "status", "data"))
			pod.SetLabels(map[string]string{"foo": "bar"})
			return *pod
		}(),
	}, {
		name: "different",
		expected: func() unstructured.Unstructured {
			pod := pod.DeepCopy()
			pod.SetResourceVersion("123")
			return *pod
		}(),
		actual: pod,
		want:   "--- expected\n+++ actual\n@@ -5,7 +5,6 @@\n     bar: baz\n   name: foo\n   namespace: bar\n-  resourceVersion: \"123\"\n spec:\n   data:\n   - foo\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrettyDiff(tt.expected, tt.actual)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
