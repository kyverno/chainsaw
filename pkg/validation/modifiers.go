package validation

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateModifiers(path *field.Path, objs ...v1alpha1.Modifier) field.ErrorList {
	var errs field.ErrorList
	for i, modifier := range objs {
		errs = append(errs, ValidateModifier(path.Index(i), modifier)...)
	}
	return errs
}
