package expressions

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

func StringPointer(ctx context.Context, in *string, bindings binding.Bindings) (*string, error) {
	if in == nil {
		return nil, nil
	}
	if *in == "" {
		return in, nil
	}
	expression := Parse(ctx, *in)
	if expression == nil || expression.Engine == "" {
		return in, nil
	}
	if converted, err := template.Execute(ctx, expression.Statement, nil, bindings, template.WithFunctionCaller(functions.Caller())); err != nil {
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
