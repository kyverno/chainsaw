package expressions

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

func StringPointer(ctx context.Context, c compilers.Compilers, in *string, bindings apis.Bindings) (*string, error) {
	if in == nil {
		return nil, nil
	}
	if *in == "" {
		return in, nil
	}
	expression := Parse(ctx, *in)
	if expression == nil {
		return in, nil
	}
	if compiler := c.Compiler(expression.Engine); compiler == nil {
		return &expression.Statement, nil
	} else if converted, err := compilers.Execute(expression.Statement, nil, bindings, compiler); err != nil {
		return nil, err
	} else if converted == nil {
		return nil, nil
	} else {
		if converted, ok := converted.(*string); ok {
			return converted, nil
		}
		if converted, ok := converted.(string); ok {
			return &converted, nil
		}
		return nil, fmt.Errorf("expression didn't evaluate to a string pointer (%s)", *in)
	}
}
