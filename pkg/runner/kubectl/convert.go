package kubectl

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

func convertString(in string, bindings binding.Bindings) (string, error) {
	if in == "" {
		return "", nil
	}
	ctx := context.TODO()
	if converted, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, in), nil, bindings, template.WithFunctionCaller(functions.Caller)); err != nil {
		return "", err
	} else {
		if covnerted, ok := converted.(string); !ok {
			return "", fmt.Errorf("expression didn't evaluate to a string (%s)", in)
		} else {
			return covnerted, nil
		}
	}
}
