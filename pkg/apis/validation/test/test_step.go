package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTestStep(path *field.Path, obj v1alpha1.TestStep) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, ValidateTestStepSpec(path, obj.TestStepSpec)...)
	return errs
}
