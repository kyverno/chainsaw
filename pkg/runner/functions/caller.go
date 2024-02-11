package functions

import (
	"context"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

var Caller = func() interpreter.FunctionCaller {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, template.GetFunctions(context.Background())...)
	funcs = append(funcs, GetFunctions()...)
	return interpreter.NewFunctionCaller(funcs...)
}()
