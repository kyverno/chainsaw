package outputs

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
)

type Outputs = map[string]any

func Process(ctx context.Context, tc binding.Bindings, input any, outputs ...v1alpha1.Output) (Outputs, error) {
	var results Outputs
	for _, output := range outputs {
		if output.Match != nil && output.Match.Value != nil {
			if errs, err := checks.Check(ctx, input, tc, output.Match); err != nil {
				return nil, err
			} else if len(errs) != 0 {
				continue
			}
		}
		name, value, err := bindings.ResolveBinding(ctx, tc, input, output.Binding)
		if err != nil {
			return nil, err
		}
		tc = bindings.RegisterBinding(ctx, tc, name, value)
		if results == nil {
			results = Outputs{}
		}
		results[name] = value
	}
	return results, nil
}
