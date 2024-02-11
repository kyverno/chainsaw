package validation

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateModifier(path *field.Path, obj v1alpha1.Modifier) field.ErrorList {
	var errs field.ErrorList
	errs = append(errs, ValidateCheck(path.Child("match"), obj.Match)...)
	count := 0
	if obj.Annotate != nil {
		count++
	}
	if obj.Label != nil {
		count++
	}
	if obj.Merge != nil {
		count++
	}
	if count == 0 {
		errs = append(errs, field.Invalid(path, obj, "no statement found in modifier"))
	} else if count > 1 {
		errs = append(errs, field.Invalid(path, obj, fmt.Sprintf("only one statement is allowed per modifier (found %d)", count)))
	} else {
		errs = append(errs, ValidateCheck(path.Child("annotate"), obj.Annotate)...)
		errs = append(errs, ValidateCheck(path.Child("label"), obj.Label)...)
		errs = append(errs, ValidateCheck(path.Child("merge"), obj.Merge)...)
	}
	return errs
}
