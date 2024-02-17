package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTestStepSpec(path *field.Path, obj v1alpha1.TestStepSpec) field.ErrorList {
	var errs field.ErrorList
	for i, try := range obj.Try {
		errs = append(errs, ValidateOperation(path.Child("try").Index(i), try)...)
	}
	for i, catch := range obj.Catch {
		errs = append(errs, ValidateCatch(path.Child("catch").Index(i), catch)...)
	}
	for i, finally := range obj.Finally {
		errs = append(errs, ValidateFinally(path.Child("finally").Index(i), finally)...)
	}
	return errs
}
