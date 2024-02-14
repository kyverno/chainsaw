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
		modifiers []v1alpha1.Any
		wantErr   bool
	}{{
		name: "nil",
		in:   in,
		out:  in,
	}, {
		name:      "empty",
		in:        in,
		modifiers: []v1alpha1.Any{},
		out:       in,
	}, {
		name: "merge",
		in:   in,
		modifiers: []v1alpha1.Any{{
			Value: map[string]any{
				"foo": "baz",
			},
		}},
		out: unstructured.Unstructured{
			Object: map[string]any{
				"foo": "baz",
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
