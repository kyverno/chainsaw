package model

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

func Evaluate[T any](ctx context.Context, in string, bindings binding.Bindings) (T, error) {
	var def T
	if converted, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, in), nil, bindings, template.WithFunctionCaller(functions.Caller)); err != nil {
		return def, err
	} else if converted, ok := converted.(T); !ok {
		return converted, fmt.Errorf("expression didn't evaluate to the expected type (%s)", in)
	} else {
		return converted, nil
	}
}

func String(ctx context.Context, in string, bindings binding.Bindings) (string, error) {
	return Evaluate[string](ctx, in, bindings)
}
