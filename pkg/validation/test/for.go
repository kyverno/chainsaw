package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFor(path *field.Path, obj *v1alpha1.For) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Deletion == nil && obj.Condition == nil {
			errs = append(errs, field.Invalid(path, obj, "either a deletion or a condition must be specified"))
		}
		if obj.Deletion != nil && obj.Condition != nil {
			errs = append(errs, field.Invalid(path, obj, "a deletion or a condition must be specified (found both)"))
		}
		if obj.Condition != nil && obj.Condition.Name == "" {
			errs = append(errs, field.Invalid(path.Child("name"), obj, "a condition name must be specified"))
		}
	}
	return errs
}
