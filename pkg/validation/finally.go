package validation

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateFinally(path *field.Path, obj v1alpha1.Finally) field.ErrorList {
	var errs field.ErrorList
	count := 0
	if obj.PodLogs != nil {
		count++
	}
	if obj.Events != nil {
		count++
	}
	if obj.Describe != nil {
		count++
	}
	if obj.Get != nil {
		count++
	}
	if obj.Command != nil {
		count++
	}
	if obj.Script != nil {
		count++
	}
	if obj.Sleep != nil {
		count++
	}
	if count == 0 {
		errs = append(errs, field.Invalid(path, obj, "no statement found in operation"))
	} else if count > 1 {
		errs = append(errs, field.Invalid(path, obj, fmt.Sprintf("only one statement is allowed per operation (found %d)", count)))
	} else {
		errs = append(errs, ValidatePodLogs(path.Child("podLogs"), obj.PodLogs)...)
		errs = append(errs, ValidateEvents(path.Child("events"), obj.Events)...)
		errs = append(errs, ValidateCommand(path.Child("command"), obj.Command)...)
		errs = append(errs, ValidateScript(path.Child("script"), obj.Script)...)
		errs = append(errs, ValidateDescribe(path.Child("describe"), obj.Describe)...)
		errs = append(errs, ValidateGet(path.Child("get"), obj.Get)...)
	}
	return errs
}
