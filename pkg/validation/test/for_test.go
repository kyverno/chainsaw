package test

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateFor(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		obj  *v1alpha1.For
		want field.ErrorList
	}{{
		name: "null",
	}, {
		name: "empty",
		path: field.NewPath("for"),
		obj:  &v1alpha1.For{},
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "for",
				BadValue: &v1alpha1.For{},
				Detail:   "either a deletion or a condition must be specified",
			},
		},
	}, {
		name: "no condition name",
		path: field.NewPath("for"),
		obj: &v1alpha1.For{
			Condition: &v1alpha1.Condition{},
		},
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "for.condition.name",
				BadValue: &v1alpha1.For{Condition: &v1alpha1.Condition{}},
				Detail:   "a condition name must be specified",
			},
		},
	}, {
		name: "both condition and deletion",
		path: field.NewPath("for"),
		obj: &v1alpha1.For{
			Condition: &v1alpha1.Condition{
				Name: "foo",
			},
			Deletion: &v1alpha1.Deletion{},
		},
		want: field.ErrorList{
			&field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "for",
				BadValue: &v1alpha1.For{
					Condition: &v1alpha1.Condition{
						Name: "foo",
					},
					Deletion: &v1alpha1.Deletion{},
				},
				Detail: "a deletion or a condition must be specified (found both)",
			},
		},
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFor(tt.path, tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}
