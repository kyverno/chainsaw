package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTest(obj *v1alpha1.Test) field.ErrorList {
	var errs field.ErrorList
	var path *field.Path
	errs = append(errs, ValidateTestSpec(path.Child("spec"), obj.Spec)...)
	return errs
}
