package test

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateTestSpec(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		obj  v1alpha1.TestSpec
		want field.ErrorList
	}{{
		name: "empty",
	}, {
		name: "invalid catch",
		obj: v1alpha1.TestSpec{
			Catch: []v1alpha1.Catch{{}},
		},
		want: field.ErrorList{
			field.Invalid(field.NewPath("catch").Index(0), v1alpha1.Catch{}, "no statement found in operation"),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateTestSpec(tt.path, tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}
