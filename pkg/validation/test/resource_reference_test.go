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
		obj  v1alpha1.ResourceReference
		want field.ErrorList
	}{{
		name: "empty",
		path: field.NewPath("foo"),
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "foo",
				BadValue: v1alpha1.ResourceReference{},
				Detail:   "kind or resource must be specified",
			},
		},
	}, {
		name: "both kind and resource",
		path: field.NewPath("foo"),
		obj: v1alpha1.ResourceReference{
			Resource: "foo",
			Kind:     "bar",
		},
		want: field.ErrorList{
			&field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "foo",
				BadValue: v1alpha1.ResourceReference{
					Resource: "foo",
					Kind:     "bar",
				},
				Detail: "kind or resource must be specified (found both)",
			},
		},
	}, {
		name: "resource and apiVersion",
		path: field.NewPath("foo"),
		obj: v1alpha1.ResourceReference{
			APIVersion: "v1",
			Resource:   "foo",
		},
		want: field.ErrorList{
			&field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "foo.apiVersion",
				BadValue: v1alpha1.ResourceReference{
					APIVersion: "v1",
					Resource:   "foo",
				},
				Detail: "apiVersion must not be specified when resource is set",
			},
		},
	}, {
		name: "kind and no apiVersion",
		path: field.NewPath("foo"),
		obj: v1alpha1.ResourceReference{
			Kind: "foo",
		},
		want: field.ErrorList{
			&field.Error{
				Type:  field.ErrorTypeInvalid,
				Field: "foo.apiVersion",
				BadValue: v1alpha1.ResourceReference{
					Kind: "foo",
				},
				Detail: "apiVersion must be specified when kind is set",
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
