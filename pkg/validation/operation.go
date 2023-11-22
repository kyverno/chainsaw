package validation

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateOperation(path *field.Path, obj v1alpha1.Operation) field.ErrorList {
	var errs field.ErrorList
	count := 0
	if obj.Apply != nil {
		count++
	}
	if obj.Assert != nil {
		count++
	}
	if obj.Command != nil {
		count++
	}
	if obj.Create != nil {
		count++
	}
	if obj.Delete != nil {
		count++
	}
	if obj.Error != nil {
		count++
	}
	if obj.Script != nil {
		count++
	}
	if count == 0 {
		errs = append(errs, field.Invalid(path, obj, "no statement found in operation"))
	} else if count > 1 {
		errs = append(errs, field.Invalid(path, obj, fmt.Sprintf("only one statement is allowed per operation (found %d)", count)))
	}
	return errs
}
