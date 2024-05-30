package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateDelete(path *field.Path, obj *v1alpha1.Delete) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		errs = append(errs, ValidateObjectReference(path.Child("ref"), obj.Ref)...)
		errs = append(errs, ValidateExpectations(path.Child("expect"), obj.Expect...)...)
		errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
	}
	return errs
}
