package check

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/check/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var caller = func() interpreter.FunctionCaller {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, template.GetFunctions(context.Background())...)
	funcs = append(funcs, functions.GetFunctions()...)
	return interpreter.NewFunctionCaller(funcs...)
}()

func Check(ctx context.Context, obj any, bindings binding.Bindings, check *v1alpha1.Check) (field.ErrorList, error) {
	if check == nil {
		return nil, errors.New("check is null")
	}
	if check.Value == nil {
		return nil, errors.New("check value is null")
	}
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return assert.Assert(ctx, nil, assert.Parse(ctx, check.Value), obj, bindings, template.WithFunctionCaller(caller))
}
