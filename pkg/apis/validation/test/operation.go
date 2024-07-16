package test

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateOperation(path *field.Path, obj v1alpha1.Operation) field.ErrorList {
	actions := [...]bool{
		obj.Apply != nil,
		obj.Assert != nil,
		obj.Command != nil,
		obj.Create != nil,
		obj.Delete != nil,
		obj.Describe != nil,
		obj.Error != nil,
		obj.Events != nil,
		obj.Get != nil,
		obj.Patch != nil,
		obj.PodLogs != nil,
		obj.Proxy != nil,
		obj.Script != nil,
		obj.Sleep != nil,
		obj.Update != nil,
		obj.Wait != nil,
	}
	count := 0
	for _, action := range actions {
		if action {
			count++
		}
	}
	var errs field.ErrorList
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
		errs = append(errs, ValidateDescribe(path.Child("describe"), obj.Describe)...)
		errs = append(errs, ValidateError(path.Child("error"), obj.Error)...)
		errs = append(errs, ValidateEvents(path.Child("events"), obj.Events)...)
		errs = append(errs, ValidateGet(path.Child("get"), obj.Get)...)
		errs = append(errs, ValidatePatch(path.Child("patch"), obj.Patch)...)
		errs = append(errs, ValidatePodLogs(path.Child("podLogs"), obj.PodLogs)...)
		errs = append(errs, ValidateProxy(path.Child("proxy"), obj.Proxy)...)
		errs = append(errs, ValidateScript(path.Child("script"), obj.Script)...)
		errs = append(errs, ValidateUpdate(path.Child("update"), obj.Update)...)
		errs = append(errs, ValidateWait(path.Child("wait"), obj.Wait)...)
	}
	return errs
}
