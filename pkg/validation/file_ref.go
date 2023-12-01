package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	util "github.com/kyverno/chainsaw/pkg/utils/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFileRef(path *field.Path, obj v1alpha1.FileRef) field.ErrorList {
	var errs field.ErrorList
	if obj.File == "" {
		errs = append(errs, field.Invalid(path.Child("file"), obj, "a file reference must be specified"))
	} else if !util.IsValidPathOrURI(obj.File) {
		errs = append(errs, field.Invalid(path.Child("file"), obj.File, "the file reference must be a valid file path or URI"))
	}
	return errs
}
