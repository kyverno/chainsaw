package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTestStep(obj v1alpha1.TestStep) field.ErrorList {
	var errs field.ErrorList
	var path *field.Path
	errs = append(errs, ValidateTestStepSpec(path.Child("spec"), obj.Spec)...)
	return errs
}
