package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTestSpec(path *field.Path, obj v1alpha1.TestSpec) field.ErrorList {
	var errs field.ErrorList
	for i, step := range obj.Steps {
		errs = append(errs, ValidateTestStep(path.Child("steps").Index(i), step)...)
	}
	for i, catch := range obj.Catch {
		errs = append(errs, ValidateCatchFinally(path.Child("catch").Index(i), catch)...)
	}
	errs = append(errs, ValidateCheck(path.Child("namespaceTemplate"), obj.NamespaceTemplate)...)
	errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
	return errs
}
