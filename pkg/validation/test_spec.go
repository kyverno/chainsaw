package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateTestSpec(path *field.Path, obj v1alpha1.TestSpec) field.ErrorList {
	var errs field.ErrorList
	for i, step := range obj.Steps {
		errs = append(errs, ValidateTestSpecStep(path.Child("steps").Index(i), step)...)
	}
	errs = append(errs, ValidateCheck(path.Child("namespaceTemplate"), obj.NamespaceTemplate)...)
	return errs
}
