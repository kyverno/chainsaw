package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFileRefOrCheck(path *field.Path, obj v1alpha1.ActionCheckRef) field.ErrorList {
	var errs field.ErrorList
	if obj.File == "" && obj.Check == nil {
		errs = append(errs, field.Invalid(path, obj, "a file reference or raw check must be specified"))
	} else if obj.File != "" && obj.Check != nil {
		errs = append(errs, field.Invalid(path, obj, "a file reference or raw check must be specified (found both)"))
	}
	return errs
}
