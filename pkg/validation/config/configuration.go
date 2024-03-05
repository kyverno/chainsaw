package config

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateConfiguration(obj *v1alpha1.Configuration) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		var path *field.Path
		errs = append(errs, ValidateConfigurationSpec(path.Child("spec"), obj.Spec)...)
	}
	return errs
}
