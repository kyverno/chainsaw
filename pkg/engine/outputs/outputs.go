package outputs

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/check"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
)

func ProcessOutputs(ctx context.Context, tc binding.Bindings, input any, outputs ...v1alpha1.Output) (operations.Outputs, error) {
	var results operations.Outputs
	for _, output := range outputs {
		if output.Match != nil && output.Match.Value != nil {
			if errs, err := check.Check(ctx, input, tc, output.Match); err != nil {
				return nil, err
			} else if len(errs) != 0 {
				continue
			}
		}
		name, value, err := bindings.ResolveBinding(ctx, tc, input, output.Binding)
		if err != nil {
			return nil, err
		}
		tc = bindings.RegisterNamedBinding(ctx, tc, name, value)
		if results == nil {
			results = operations.Outputs{}
		}
		results[name] = value
	}
	return results, nil
}
