package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateCreate(path *field.Path, obj *v1alpha1.Create) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		errs = append(errs, ValidateFileRefOrResource(path, obj.ActionResourceRef)...)
		errs = append(errs, ValidateExpectations(path.Child("expect"), obj.Expect...)...)
		errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
	}
	return errs
}
