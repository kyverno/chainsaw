package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateResourceReference(path *field.Path, obj v1alpha1.ResourceReference) field.ErrorList {
	var errs field.ErrorList
	if obj.Kind == "" && obj.Resource == "" {
		errs = append(errs, field.Invalid(path, obj, "kind or resource must be specified"))
	} else if obj.Kind != "" && obj.Resource != "" {
		errs = append(errs, field.Invalid(path, obj, "kind or resource must be specified (found both)"))
	} else if obj.Resource != "" && obj.APIVersion != "" {
		errs = append(errs, field.Invalid(path.Child("apiVersion"), obj, "apiVersion must not be specified when resource is set"))
	} else if obj.Kind != "" && obj.APIVersion == "" {
		errs = append(errs, field.Invalid(path.Child("apiVersion"), obj, "apiVersion must be specified when kind is set"))
	}
	return errs
}
