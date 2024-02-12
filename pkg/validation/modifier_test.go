package validation

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateModifier(t *testing.T) {
	multiple := v1alpha1.Modifier{
		Label: &v1alpha1.Any{
			Value: map[string]any{
				"foo": "bar",
			},
		},
		Merge: &v1alpha1.Any{
			Value: map[string]any{
				"foo": "bar",
			},
		},
	}
	tests := []struct {
		name     string
		path     *field.Path
		modifier v1alpha1.Modifier
		want     field.ErrorList
	}{{
		name: "invalid match",
		path: field.NewPath("foo"),
		modifier: v1alpha1.Modifier{
			Match: &v1alpha1.Match{},
			Label: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo").Child("match"), &v1alpha1.Check{}, "a value must be specified"),
		},
	}, {
		name:     "no statement",
		path:     field.NewPath("foo"),
		modifier: v1alpha1.Modifier{},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo"), v1alpha1.Modifier{}, "no statement found in modifier"),
		},
	}, {
		name: "annotate",
		path: field.NewPath("foo"),
		modifier: v1alpha1.Modifier{
			Annotate: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		},
		want: nil,
	}, {
		name: "label",
		path: field.NewPath("foo"),
		modifier: v1alpha1.Modifier{
			Label: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		},
		want: nil,
	}, {
		name: "merge",
		path: field.NewPath("foo"),
		modifier: v1alpha1.Modifier{
			Merge: &v1alpha1.Any{
				Value: map[string]any{
					"foo": "bar",
				},
			},
		},
		want: nil,
	}, {
		name:     "multiple statements",
		path:     field.NewPath("foo"),
		modifier: multiple,
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo"), multiple, "only one statement is allowed per modifier (found 2)"),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateModifier(tt.path, tt.modifier)
			assert.Equal(t, tt.want, got)
		})
	}
}
