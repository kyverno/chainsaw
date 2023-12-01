package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	util "github.com/kyverno/chainsaw/pkg/utils/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFileRefOrResource(path *field.Path, obj v1alpha1.FileRefOrResource) field.ErrorList {
	var errs field.ErrorList
	if obj.File == "" && obj.Resource == nil {
		errs = append(errs, field.Invalid(path, obj, "a file reference or raw resource must be specified"))
	} else if obj.File != "" && obj.Resource != nil {
		errs = append(errs, field.Invalid(path, obj, "a file reference or raw resource must be specified (found both)"))
	} else if obj.File != "" && !util.IsValidPathOrURI(obj.File) {
		errs = append(errs, field.Invalid(path.Child("file"), obj.File, "the file reference must be a valid file path or URI"))
	}
	return errs
}
