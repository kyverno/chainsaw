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
			Object: map[string]any{
				"": func() {},
			},
		},
		wantErr: true,
	}, {
		name: "error",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"": "",
			},
		},
		actual: unstructured.Unstructured{
			Object: map[string]any{
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
	}, {
		// expected holds a scalar where actual holds a map for the same key.
		// pruning must not assume the expected value is also a map.
		name: "type mismatch scalar vs map",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"data":       "scalar",
			},
		},
		actual: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"data":       map[string]any{"key": "value"},
			},
		},
		want: "--- expected\n+++ actual\n@@ -1,4 +1,5 @@\n apiVersion: v1\n-data: scalar\n+data:\n+  key: value\n kind: ConfigMap\n \n",
	}, {
		// same mismatch nested under metadata, which uses a dedicated prune path.
		name: "type mismatch under metadata",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]any{
					"name":   "foo",
					"labels": "scalar",
				},
			},
		},
		actual: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]any{
					"name":   "foo",
					"labels": map[string]any{"app": "bar"},
				},
			},
		},
		want: "--- expected\n+++ actual\n@@ -1,6 +1,7 @@\n apiVersion: v1\n kind: ConfigMap\n metadata:\n-  labels: scalar\n+  labels:\n+    app: bar\n   name: foo\n \n",
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
