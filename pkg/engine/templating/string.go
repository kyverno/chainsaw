package templating

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func String(ctx context.Context, in string, bindings binding.Bindings) (string, error) {
	if in == "" {
		return "", nil
	}
	if converted, err := Template(ctx, v1alpha1.Any{Value: in}, nil, bindings); err != nil {
		return "", err
	} else {
		if converted, ok := converted.(string); !ok {
			return "", fmt.Errorf("expression didn't evaluate to a string (%s)", in)
		} else {
			return converted, nil
		}
	}
}
