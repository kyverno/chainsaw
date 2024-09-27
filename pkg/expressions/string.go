package expressions

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
)

func String(ctx context.Context, in string, bindings binding.Bindings) (string, error) {
	if in == "" {
		return "", nil
	}
	expression := Parse(ctx, in)
	if expression == nil || expression.Engine == "" {
		return in, nil
	}
	if converted, err := compilers.Execute(expression.Statement, nil, bindings, jp.NewCompiler(jp.WithFunctionCaller(functions.Caller()))); err != nil {
		return "", err
	} else {
		if converted, ok := converted.(string); !ok {
			return "", fmt.Errorf("expression didn't evaluate to a string (%s)", in)
		} else {
			return converted, nil
		}
	}
}
