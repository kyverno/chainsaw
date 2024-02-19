package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateWait(path *field.Path, obj *v1alpha1.Wait) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Resource == "" {
			errs = append(errs, field.Invalid(path.Child("resource"), obj, "a resource must be specified"))
		}
		if obj.Name != "" && obj.Selector != "" {
			errs = append(errs, field.Invalid(path, obj, "a name or label selector must be specified (found both)"))
		}
	}
	return errs
}
