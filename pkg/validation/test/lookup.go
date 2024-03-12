package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateLookup(path *field.Path, obj *v1alpha1.Lookup) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		errs = append(errs, ValidateFileRefOrResource(path, obj.FileRefOrResource)...)
		errs = append(errs, ValidateBindings(path.Child("bindings"), obj.Bindings...)...)
		errs = append(errs, ValidateOutputs(path.Child("outputs"), obj.Outputs...)...)
	}
	return errs
}
