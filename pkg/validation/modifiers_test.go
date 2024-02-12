package validation

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateModifiers(t *testing.T) {
	tests := []struct {
		name      string
		path      *field.Path
		modifiers []v1alpha1.Modifier
		want      field.ErrorList
	}{{
		name:      "nil",
		path:      nil,
		modifiers: nil,
		want:      nil,
	}, {
		name:      "empty",
		path:      field.NewPath("foo"),
		modifiers: []v1alpha1.Modifier{},
		want:      nil,
	}, {
		name: "invalid match",
		path: field.NewPath("foo"),
		modifiers: []v1alpha1.Modifier{{
			Match: &v1alpha1.Match{},
			Label: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		}},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo").Index(0).Child("match"), &v1alpha1.Check{}, "a value must be specified"),
		},
	}, {
		name: "valid",
		path: field.NewPath("foo"),
		modifiers: []v1alpha1.Modifier{{
			Annotate: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		}},
		want: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateModifiers(tt.path, tt.modifiers...)
			assert.Equal(t, tt.want, got)
		})
	}
}
