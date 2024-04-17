package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFor(path *field.Path, obj *v1alpha1.For) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		if obj.Deletion == nil && obj.Condition == nil && obj.JsonPath == nil {
			errs = append(errs, field.Invalid(path, obj, "either a deletion, condition or a jsonpath must be specified"))
		}
		if (obj.Deletion != nil && obj.Condition != nil) || (obj.Deletion != nil && obj.JsonPath != nil) || (obj.Condition != nil && obj.JsonPath != nil) {
			errs = append(errs, field.Invalid(path, obj, "a deletion, condition or jsonpath must be specified (only one)"))
		}
		if obj.Condition != nil && obj.Condition.Name == "" {
			errs = append(errs, field.Invalid(path.Child("condition").Child("name"), obj, "a condition name must be specified"))
		}
		if obj.JsonPath != nil && obj.JsonPath.Path == "" {
			errs = append(errs, field.Invalid(path.Child("jsonPath").Child("path"), obj, "a json path must be specified"))
		}
		if obj.JsonPath != nil && obj.JsonPath.Value == "" {
			errs = append(errs, field.Invalid(path.Child("jsonPath").Child("value"), obj, "a value must be specified"))
		}
	}
	return errs
}
