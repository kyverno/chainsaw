package checks

import (
	"context"
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Check(ctx context.Context, compilers compilers.Compilers, obj any, bindings apis.Bindings, check *v1alpha1.Check) (field.ErrorList, error) {
	if check == nil {
		return nil, errors.New("check is null")
	}
	if check.IsNil() {
		return nil, errors.New("check value is null")
	}
	if assertion, err := check.Compile(nil, compilers); err != nil {
		return nil, err
	} else {
		if bindings == nil {
			bindings = apis.NewBindings()
		}
		return assertion.Assert(nil, obj, bindings)
	}
}
