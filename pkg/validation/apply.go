package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateApply(path *field.Path, obj *v1alpha1.Apply) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		errs = append(errs, ValidateFileRefOrResource(path, obj.FileRefOrResource)...)
		errs = append(errs, ValidateExpectations(path.Child("expect"), obj.Expect...)...)
		errs = append(errs, ValidateModifiers(path.Child("modifiers"), obj.Modifiers...)...)
	}
	return errs
}
