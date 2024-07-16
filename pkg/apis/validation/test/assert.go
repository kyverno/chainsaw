package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateAssert(path *field.Path, obj *v1alpha1.Assert) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		errs = append(errs, ValidateFileRefOrCheck(path, obj.ActionCheckRef)...)
		errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
	}
	return errs
}
