package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFileRef(path *field.Path, obj v1alpha1.FileRef) field.ErrorList {
	var errs field.ErrorList
	if obj.File == "" {
		errs = append(errs, field.Invalid(path, obj, "a file reference must be specified"))
	}
	return errs
}
