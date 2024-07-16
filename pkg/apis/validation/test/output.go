package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateOutput(path *field.Path, obj v1alpha1.Output) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, ValidateBinding(path, obj.Binding)...)
	return errs
}

func ValidateOutputs(path *field.Path, objs ...v1alpha1.Output) field.ErrorList {
	var errs field.ErrorList
	for i, obj := range objs {
		errs = append(errs, ValidateOutput(path.Index(i), obj)...)
	}
	return errs
}
