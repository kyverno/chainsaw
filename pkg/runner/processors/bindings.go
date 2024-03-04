package processors

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/client-go/rest"
)

func registerBindings(
	ctx context.Context,
	bindings binding.Bindings,
	config *rest.Config,
	client client.Client,
	variables ...v1alpha1.Binding,
) (binding.Bindings, error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	bindings = bindings.Register("$client", binding.NewBinding(client))
	bindings = bindings.Register("$config", binding.NewBinding(config))
	for _, variable := range variables {
		if err := variable.CheckName(); err != nil {
			return bindings, err
		}
		patched, err := mutation.Mutate(ctx, nil, mutation.Parse(ctx, variable.Value.Value), nil, bindings, template.WithFunctionCaller(functions.Caller))
		if err != nil {
			return bindings, err
		}
		bindings = bindings.Register("$"+variable.Name, binding.NewBinding(patched))
	}
	return bindings, nil
}
