package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFileRefOrResource(path *field.Path, obj v1alpha1.FileRefOrResource) field.ErrorList {
	var errs field.ErrorList
	if obj.File == "" && obj.Resource == nil {
		errs = append(errs, field.Invalid(path, obj, "a file reference or raw resource must be specified"))
	} else if obj.File != "" && obj.Resource != nil {
		errs = append(errs, field.Invalid(path, obj, "a file reference or raw resource must be specified (found both)"))
	}
	return errs
}
