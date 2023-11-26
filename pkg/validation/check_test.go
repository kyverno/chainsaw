package validation

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateCheck(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		obj  *v1alpha1.Check
		want field.ErrorList
	}{{
		name: "nil",
		path: nil,
		obj:  nil,
		want: nil,
	}, {
		name: "no value",
		path: field.NewPath("foo"),
		obj:  &v1alpha1.Check{},
		want: field.ErrorList{
			field.Invalid(field.NewPath("foo"), &v1alpha1.Check{}, "a value must be specified"),
		},
	}, {
		name: "valid",
		path: field.NewPath("foo"),
		obj: &v1alpha1.Check{
			Value: 42,
		},
		want: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCheck(tt.path, tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}
