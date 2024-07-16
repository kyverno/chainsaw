package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateBinding(path *field.Path, obj v1alpha1.Binding) field.ErrorList {
	var errs field.ErrorList
	if err := obj.CheckName(); err != nil {
		errs = append(errs, field.Invalid(path.Child("name"), obj.Name, err.Error()))
	}
	return errs
}

func ValidateBindings(path *field.Path, objs ...v1alpha1.Binding) field.ErrorList {
	var errs field.ErrorList
	for i, obj := range objs {
		errs = append(errs, ValidateBinding(path.Index(i), obj)...)
	}
	return errs
}
