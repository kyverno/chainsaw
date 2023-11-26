package validation

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateExpectations(t *testing.T) {
	tests := []struct {
		name         string
		path         *field.Path
		expectations []v1alpha1.Expectation
		want         field.ErrorList
	}{{
		name:         "nil",
		path:         nil,
		expectations: nil,
		want:         nil,
	}, {
		name: "invalid check",
		path: field.NewPath("foo"),
		expectations: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{},
		}},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo").Index(0).Child("check"), &v1alpha1.Check{}, "a value must be specified"),
		},
	}, {
		name: "invalid match",
		path: field.NewPath("foo"),
		expectations: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{},
			Check: v1alpha1.Check{
				Value: 42,
			},
		}},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo").Index(0).Child("match"), &v1alpha1.Check{}, "a value must be specified"),
		},
	}, {
		name: "invalid check and match",
		path: field.NewPath("foo"),
		expectations: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{},
			Check: v1alpha1.Check{},
		}},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo").Index(0).Child("match"), &v1alpha1.Check{}, "a value must be specified"),
			field.Invalid(field.NewPath("foo").Index(0).Child("check"), &v1alpha1.Check{}, "a value must be specified"),
		},
	}, {
		name: "valid",
		path: field.NewPath("foo"),
		expectations: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: 42,
			},
			Check: v1alpha1.Check{
				Value: 42,
			},
		}},
		want: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateExpectations(tt.path, tt.expectations...)
			assert.Equal(t, tt.want, got)
		})
	}
}
