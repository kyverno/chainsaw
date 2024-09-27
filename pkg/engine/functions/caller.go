package functions

import (
	"context"
	"sync"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/kyverno/kyverno-json/pkg/jp"
)

var Caller = sync.OnceValue(func() interpreter.FunctionCaller {
	var funcs []jpfunctions.FunctionEntry
	funcs = append(funcs, jp.GetFunctions(context.Background())...)
	funcs = append(funcs, GetFunctions()...)
	return interpreter.NewFunctionCaller(funcs...)
})
