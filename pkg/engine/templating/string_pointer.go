package templating

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func StringPointer(ctx context.Context, in *string, bindings binding.Bindings) (*string, error) {
	if in == nil {
		return nil, nil
	}
	if *in == "" {
		return in, nil
	}
	if converted, err := Template(ctx, v1alpha1.Any{Value: *in}, nil, bindings); err != nil {
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
