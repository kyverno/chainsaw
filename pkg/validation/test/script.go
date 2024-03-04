package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateScript(path *field.Path, obj *v1alpha1.Script) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Content == "" {
			errs = append(errs, field.Invalid(path.Child("content"), obj, "content must be specified"))
		}
		errs = append(errs, ValidateCheck(path.Child("check"), obj.Check)...)
	}
	return errs
}
