package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateCommand(path *field.Path, obj *v1alpha1.Command) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Entrypoint == "" {
			errs = append(errs, field.Invalid(path.Child("entrypoint"), obj, "entrypoint must be specified"))
		}
		errs = append(errs, ValidateCheck(path.Child("check"), obj.Check)...)
		errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
	}
	return errs
}
