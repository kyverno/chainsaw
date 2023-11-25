package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateCheck(path *field.Path, obj *v1alpha1.Check) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Value == nil {
			errs = append(errs, field.Invalid(path, obj, "a value must be specified"))
		}
	}
	return errs
}
