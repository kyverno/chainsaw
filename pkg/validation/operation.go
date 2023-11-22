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
	} else {
		if obj.Apply != nil {
			errs = append(errs, ValidateFileRefOrResource(path.Child("apply"), obj.Apply.FileRefOrResource)...)
		}
		if obj.Assert != nil {
			errs = append(errs, ValidateFileRef(path.Child("assert"), obj.Assert.FileRef)...)
		}
		// TODO
		// if obj.Command != nil {
		// }
		if obj.Create != nil {
			errs = append(errs, ValidateFileRefOrResource(path.Child("create"), obj.Create.FileRefOrResource)...)
		}
		// TODO
		// 	if obj.Delete != nil {
		// }
		if obj.Error != nil {
			errs = append(errs, ValidateFileRef(path.Child("error"), obj.Error.FileRef)...)
		}
		// TODO
		// 	if obj.Script != nil {
		// }
	}
	return errs
}
