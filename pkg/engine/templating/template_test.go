package templating

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestTemplateAndMerge(t *testing.T) {
	in := unstructured.Unstructured{
		Object: map[string]any{
			"foo": "bar",
		},
	}
	tests := []struct {
		name      string
		in        unstructured.Unstructured
		out       unstructured.Unstructured
		templates []v1alpha1.Projection
		wantErr   bool
	}{{
		name: "nil",
		in:   in,
		out:  in,
	}, {
		name:      "empty",
		in:        in,
		templates: []v1alpha1.Projection{},
		out:       in,
	}, {
		name: "merge",
		in:   in,
		templates: []v1alpha1.Projection{
			v1alpha1.NewProjection(
				map[string]any{
					"foo": "baz",
				},
			),
		},
		out: unstructured.Unstructured{
			Object: map[string]any{
				"foo": "baz",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemplateAndMerge(context.TODO(), apis.XDefaultCompilers, tt.in, nil, tt.templates...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.out, got)
			}
		})
	}
}
