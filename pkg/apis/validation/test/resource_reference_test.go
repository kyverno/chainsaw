package test

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateResourceReference(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		obj  v1alpha1.ObjectType
		want field.ErrorList
	}{{
		name: "empty",
		path: field.NewPath("foo"),
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "foo.kind",
				BadValue: v1alpha1.ObjectType{},
				Detail:   "kind must be specified",
			},
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "foo.apiVersion",
				BadValue: v1alpha1.ObjectType{},
				Detail:   "apiVersion must be specified",
			},
		},
	}, {
		name: "kind and no apiVersion",
		path: field.NewPath("foo"),
		obj: v1alpha1.ObjectType{
			Kind: "foo",
		},
		want: field.ErrorList{
			&field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "foo.apiVersion",
				BadValue: v1alpha1.ObjectType{
					Kind: "foo",
				},
				Detail: "apiVersion must be specified",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateResourceReference(tt.path, tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}
