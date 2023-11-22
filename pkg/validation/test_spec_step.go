package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTestSpecStep(path *field.Path, obj v1alpha1.TestSpecStep) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, ValidateTestStepSpec(path, obj.Spec)...)
	return errs
}
