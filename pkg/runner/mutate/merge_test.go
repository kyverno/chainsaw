package mutate

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestMerge(t *testing.T) {
	in := unstructured.Unstructured{
		Object: map[string]any{
			"foo": "bar",
		},
	}
	tests := []struct {
		name      string
		in        unstructured.Unstructured
		out       unstructured.Unstructured
		modifiers []v1alpha1.Modifier
		wantErr   bool
	}{{
		name: "nil",
		in:   in,
		out:  in,
	}, {
		name:      "empty",
		in:        in,
		modifiers: []v1alpha1.Modifier{},
		out:       in,
	}, {
		name: "annotate",
		in:   in,
		modifiers: []v1alpha1.Modifier{{
			Annotate: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		}},
		out: unstructured.Unstructured{
			Object: map[string]any{
				"foo": "bar",
				"metadata": map[string]any{
					"annotations": map[string]any{
						"foo": "bar",
					},
				},
			},
		},
	}, {
		name: "label",
		in:   in,
		modifiers: []v1alpha1.Modifier{{
			Label: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		}},
		out: unstructured.Unstructured{
			Object: map[string]any{
				"foo": "bar",
				"metadata": map[string]any{
					"labels": map[string]any{
						"foo": "bar",
					},
				},
			},
		},
	}, {
		name: "merge",
		in:   in,
		modifiers: []v1alpha1.Modifier{{
			Merge: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "baz",
				},
			},
		}},
		out: unstructured.Unstructured{
			Object: map[string]any{
				"foo": "baz",
			},
		},
	}, {
		name: "merge",
		in:   in,
		modifiers: []v1alpha1.Modifier{{
			Match: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "baz",
				},
			},
			Merge: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		}},
		out: unstructured.Unstructured{
			Object: map[string]any{
				"foo": "bar",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Merge(context.TODO(), tt.in, nil, tt.modifiers...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out, got)
			}
		})
	}
}
