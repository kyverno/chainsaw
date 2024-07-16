package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateProxy(path *field.Path, obj *v1alpha1.Proxy) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Name == "" {
			errs = append(errs, field.Invalid(path, obj, "name must be specified"))
		}
		if obj.Namespace == "" {
			errs = append(errs, field.Invalid(path, obj, "namespace must be specified"))
		}
		errs = append(errs, ValidateResourceReference(path, obj.ObjectType)...)
	}
	return errs
}
