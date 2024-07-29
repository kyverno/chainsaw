package templating

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

func Template(ctx context.Context, tpl v1alpha1.Any, value any, bindings binding.Bindings) (any, error) {
	return mutate.Mutate(ctx, nil, mutate.Parse(ctx, tpl.Value), value, bindings, template.WithFunctionCaller(functions.Caller()))
}
