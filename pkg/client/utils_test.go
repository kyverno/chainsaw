package client

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	truntime "k8s.io/apimachinery/pkg/runtime/testing"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestObjectKey(t *testing.T) {
	obj := &metav1.ObjectMeta{
		Name:      "test-name",
		Namespace: "test-namespace",
	}

	key := ObjectKey(obj)
	assert.Equal(t, "test-name", key.Name)
	assert.Equal(t, "test-namespace", key.Namespace)
}

func TestName(t *testing.T) {
	key := ctrlclient.ObjectKey{Name: "test-name"}
	name := Name(key)
	assert.Equal(t, "test-name", name)

	key.Namespace = "test-namespace"
	name = Name(key)
	assert.Equal(t, "test-namespace/test-name", name)

	key = ctrlclient.ObjectKey{}
	name = Name(key)
	assert.Equal(t, "*", name)
}

func TestColouredName(t *testing.T) {
	disabled := color.New(color.FgBlue)
	disabled.DisableColor()
	enabled := color.New(color.FgBlue)
	enabled.EnableColor()
	tests := []struct {
		name  string
		key   ctrlclient.ObjectKey
		color *color.Color
		want  string
	}{{
		name:  "empty",
		key:   ctrlclient.ObjectKey{},
		color: nil,
		want:  "*",
	}, {
		name:  "name only",
		key:   ctrlclient.ObjectKey{Name: "test-name"},
		color: nil,
		want:  "test-name",
	}, {
		name:  "name and namespace",
		key:   ctrlclient.ObjectKey{Name: "test-name", Namespace: "test-namespace"},
		color: nil,
		want:  "test-namespace/test-name",
	}, {
		name:  "empty",
		key:   ctrlclient.ObjectKey{},
		color: disabled,
		want:  "*",
	}, {
		name:  "name only",
		key:   ctrlclient.ObjectKey{Name: "test-name"},
		color: disabled,
		want:  "test-name",
	}, {
		name:  "name and namespace",
		key:   ctrlclient.ObjectKey{Name: "test-name", Namespace: "test-namespace"},
		color: disabled,
		want:  "test-namespace/test-name",
	}, {
		name:  "empty",
		key:   ctrlclient.ObjectKey{},
		color: enabled,
		want:  "\x1b[34m*\x1b[0m",
	}, {
		name:  "name only",
		key:   ctrlclient.ObjectKey{Name: "test-name"},
		color: enabled,
		want:  "\x1b[34mtest-name\x1b[0m",
	}, {
		name:  "name and namespace",
		key:   ctrlclient.ObjectKey{Name: "test-name", Namespace: "test-namespace"},
		color: enabled,
		want:  "\x1b[34mtest-namespace\x1b[0m/\x1b[34mtest-name\x1b[0m",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColouredName(tt.key, tt.color)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPatchObject(t *testing.T) {
	tests := []struct {
		name     string
		actual   runtime.Object
		expected runtime.Object
		want     runtime.Object
		wantErr  bool
	}{{
		name:     "acutal nil",
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
