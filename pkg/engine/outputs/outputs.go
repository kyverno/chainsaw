package outputs

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

type Outputs = map[string]any

func Process(ctx context.Context, compilers compilers.Compilers, tc apis.Bindings, input any, outputs ...v1alpha1.Output) (Outputs, error) {
	var results Outputs
	for _, output := range outputs {
		if output.Match != nil && !output.Match.IsNil() {
			if errs, err := checks.Check(ctx, compilers, input, tc, output.Match); err != nil {
				return nil, err
			} else if len(errs) != 0 {
				continue
			}
		}
		name, value, err := bindings.ResolveBinding(ctx, compilers, tc, input, output.Binding)
		if err != nil {
			return nil, err
		}
		tc = bindings.RegisterBinding(tc, name, value)
		if results == nil {
			results = Outputs{}
		}
		results[name] = value
	}
	return results, nil
}
