package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateDelete(path *field.Path, obj *v1alpha1.Delete) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.File == "" && obj.Ref == nil {
			errs = append(errs, field.Invalid(path, obj, "a file or ref must be specified"))
		} else if obj.File != "" && obj.Ref != nil {
			errs = append(errs, field.Invalid(path, obj, "a file or ref must be specified (found both)"))
		} else if obj.Ref != nil {
			errs = append(errs, ValidateObjectReference(path.Child("ref"), *obj.Ref)...)
		}
		errs = append(errs, ValidateExpectations(path.Child("expect"), obj.Expect...)...)
		errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
	}
	return errs
}
