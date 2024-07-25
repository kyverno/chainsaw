package bindings

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

func String(in string, bindings binding.Bindings) (string, error) {
	if in == "" {
		return "", nil
	}
	ctx := context.TODO()
	if converted, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, in), nil, bindings, template.WithFunctionCaller(functions.Caller)); err != nil {
		return "", err
	} else {
		if converted, ok := converted.(string); !ok {
			return "", fmt.Errorf("expression didn't evaluate to a string (%s)", in)
		} else {
			return converted, nil
		}
	}
}

func StringPointer(in *string, bindings binding.Bindings) (*string, error) {
	if in == nil {
		return nil, nil
	}
	if *in == "" {
		return in, nil
	}
	ctx := context.TODO()
	if converted, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, in), nil, bindings, template.WithFunctionCaller(functions.Caller)); err != nil {
		return nil, err
	} else if converted == nil {
		return nil, nil
	} else {
		if converted, ok := converted.(*string); ok {
			return converted, nil
		}
		if converted, ok := converted.(string); !ok {
			return &converted, nil
		}
		return nil, fmt.Errorf("expression didn't evaluate to a string pointer (%s)", *in)
	}
}
