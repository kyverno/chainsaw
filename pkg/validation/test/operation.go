package test

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
	if obj.Patch != nil {
		count++
	}
	if obj.Script != nil {
		count++
	}
	if obj.Sleep != nil {
		count++
	}
	if obj.Update != nil {
		count++
	}
	if obj.Wait != nil {
		count++
	}
	if count == 0 {
		errs = append(errs, field.Invalid(path, obj, "no statement found in operation"))
	} else if count > 1 {
		errs = append(errs, field.Invalid(path, obj, fmt.Sprintf("only one statement is allowed per operation (found %d)", count)))
	} else {
		errs = append(errs, ValidateApply(path.Child("apply"), obj.Apply)...)
		errs = append(errs, ValidateAssert(path.Child("assert"), obj.Assert)...)
		errs = append(errs, ValidateCommand(path.Child("command"), obj.Command)...)
		errs = append(errs, ValidateCreate(path.Child("create"), obj.Create)...)
		errs = append(errs, ValidateDelete(path.Child("delete"), obj.Delete)...)
		errs = append(errs, ValidateError(path.Child("error"), obj.Error)...)
		errs = append(errs, ValidatePatch(path.Child("patch"), obj.Patch)...)
		errs = append(errs, ValidateScript(path.Child("script"), obj.Script)...)
		errs = append(errs, ValidateUpdate(path.Child("update"), obj.Update)...)
		errs = append(errs, ValidateWait(path.Child("wait"), obj.Wait)...)
	}
	return errs
}
