package outputs

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name    string
		tc      binding.Bindings
		input   any
		outputs []v1alpha1.Output
		want    Outputs
		wantErr bool
	}{{
		name:    "empty",
		tc:      binding.NewBindings(),
		input:   nil,
		outputs: nil,
		want:    nil,
		wantErr: false,
	}, {
		name:  "simple",
		tc:    binding.NewBindings(),
		input: nil,
		outputs: []v1alpha1.Output{{
			Binding: v1alpha1.Binding{
				Name:  "foo",
				Value: v1alpha1.NewProjection("bar"),
			},
		}},
		want: Outputs{
			"foo": "bar",
		},
		wantErr: false,
	}, {
		name:  "match",
		tc:    binding.NewBindings(),
		input: map[string]any{},
		outputs: []v1alpha1.Output{{
			Match: ptr.To(
				v1alpha1.NewMatch(
					map[string]any{
						"bar": "baz",
					},
				),
			),
			Binding: v1alpha1.Binding{
				Name:  "foo",
				Value: v1alpha1.NewProjection("bar"),
			},
		}},
		want:    nil,
		wantErr: false,
	}, {
		name:  "match error",
		tc:    binding.NewBindings(),
		input: nil,
		outputs: []v1alpha1.Output{{
			Match: ptr.To(
				v1alpha1.NewMatch(
					map[string]any{
						"($bar)": "baz",
					},
				),
			),
			Binding: v1alpha1.Binding{
				Name:  "foo",
				Value: v1alpha1.NewProjection("bar"),
			},
		}},
		want:    nil,
		wantErr: true,
	}, {
		name:  "error",
		tc:    binding.NewBindings(),
		input: nil,
		outputs: []v1alpha1.Output{{
			Binding: v1alpha1.Binding{
				Name:  "($foo)",
				Value: v1alpha1.NewProjection("bar"),
			},
		}},
		want:    nil,
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Process(context.TODO(), tt.tc, tt.input, tt.outputs...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
