package functions

import (
	"context"
	"sync"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/jmespath-community/go-jmespath/pkg/interpreter"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

var CallerFunctions = sync.OnceValue(func() []FunctionEntry {
	var funcs []FunctionEntry
	for _, function := range template.GetFunctions(context.Background()) {
		funcs = append(funcs, FunctionEntry{
			FunctionEntry: function,
		})
	}
	funcs = append(funcs, GetFunctions()...)
	return funcs
})

var Caller = sync.OnceValue(func() interpreter.FunctionCaller {
	var funcs []jpfunctions.FunctionEntry
	for _, function := range CallerFunctions() {
		funcs = append(funcs, function.FunctionEntry)
	}
	return interpreter.NewFunctionCaller(funcs...)
})
