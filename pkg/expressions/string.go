package expressions

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

func String(ctx context.Context, in string, bindings apis.Bindings) (string, error) {
	if in == "" {
		return "", nil
	}
	expression := Parse(ctx, in)
	if expression == nil {
		return in, nil
	}
	if compiler := apis.DefaultCompilers.Compiler(expression.Engine); compiler == nil {
		return expression.Statement, nil
	} else if converted, err := compilers.Execute(expression.Statement, nil, bindings, compiler); err != nil {
		return "", err
	} else {
		if converted, ok := converted.(string); !ok {
			return "", fmt.Errorf("expression didn't evaluate to a string (%s)", in)
		} else {
			return converted, nil
		}
	}
}
