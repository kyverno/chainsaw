package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateExpectations(path *field.Path, expectations ...v1alpha1.Expectation) field.ErrorList {
	var errs field.ErrorList
	for i := range expectations {
		path := path.Index(i)
		expectation := expectations[i]
		errs = append(errs, ValidateCheck(path.Child("match"), expectation.Match)...)
		errs = append(errs, ValidateCheck(path.Child("check"), &expectation.Check)...)
	}
	return errs
}
