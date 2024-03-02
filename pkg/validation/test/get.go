package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateGet(path *field.Path, obj *v1alpha1.Get) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Name != "" && obj.Selector != "" {
			errs = append(errs, field.Invalid(path, obj, "a name or label selector must be specified (found both)"))
		}
		errs = append(errs, ValidateResourceReference(path, obj.ResourceReference)...)
	}
	return errs
}
