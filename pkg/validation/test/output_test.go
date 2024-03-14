package test

import (
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateOutput(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		obj  v1alpha1.Output
		want field.ErrorList
	}{{
		name: "empty",
		path: field.NewPath("foo"),
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "foo.name",
				BadValue: "",
				Detail:   "invalid name ",
			},
		},
	}, {
		name: "valid name",
		path: field.NewPath("foo"),
		obj: v1alpha1.Output{
			Binding: v1alpha1.Binding{
				Name: "foo",
			},
		},
	}, {
		name: "invalid name",
		path: field.NewPath("foo"),
		obj: v1alpha1.Output{
			Binding: v1alpha1.Binding{
				Name: "$foo",
			},
		},
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "foo.name",
				BadValue: "$foo",
				Detail:   "invalid name $foo",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOutput(tt.path, tt.obj)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateOutputs(t *testing.T) {
	tests := []struct {
		name string
		path *field.Path
		objs []v1alpha1.Output
		want field.ErrorList
	}{{
		name: "null",
	}, {
		name: "empty",
		objs: []v1alpha1.Output{},
	}, {
		name: "valid",
		objs: []v1alpha1.Output{{
			Binding: v1alpha1.Binding{
				Name: "foo",
			},
		}},
	}, {
		name: "invalid",
		objs: []v1alpha1.Output{{}},
		want: field.ErrorList{
			&field.Error{
				Type:     field.ErrorTypeInvalid,
				Field:    "[0].name",
				BadValue: "",
				Detail:   "invalid name ",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOutputs(tt.path, tt.objs...)
			assert.Equal(t, tt.want, got)
		})
	}
}
