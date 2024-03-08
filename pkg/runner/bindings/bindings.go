package bindings

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	apitemplate "github.com/kyverno/chainsaw/pkg/runner/template"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/client-go/rest"
)

func RegisterNamedBinding(ctx context.Context, bindings binding.Bindings, name string, value any) binding.Bindings {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return bindings.Register("$"+name, binding.NewBinding(value))
}

func ResolveBinding(ctx context.Context, bindings binding.Bindings, input any, variable v1alpha1.Binding) (string, any, error) {
	name, err := apitemplate.ConvertString(variable.Name, bindings)
	if err != nil {
		return "", nil, err
	}
	if err := v1alpha1.CheckBindingName(name); err != nil {
		return "", nil, err
	}
	value, err := mutation.Mutate(ctx, nil, mutation.Parse(ctx, variable.Value.Value), input, bindings, template.WithFunctionCaller(functions.Caller))
	if err != nil {
		return "", nil, err
	}
	return name, value, err
}

func RegisterBinding(ctx context.Context, bindings binding.Bindings, input any, variable v1alpha1.Binding) (binding.Bindings, error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	name, value, err := ResolveBinding(ctx, bindings, input, variable)
	if err != nil {
		return bindings, err
	}
	return RegisterNamedBinding(ctx, bindings, name, value), nil
}

func RegisterBindings(ctx context.Context, bindings binding.Bindings, config *rest.Config, client client.Client, variables ...v1alpha1.Binding) (binding.Bindings, error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	bindings = bindings.Register("$client", binding.NewBinding(client))
	bindings = bindings.Register("$config", binding.NewBinding(config))
	for _, variable := range variables {
		next, err := RegisterBinding(ctx, bindings, nil, variable)
		if err != nil {
			return bindings, err
		}
		bindings = next
	}
	return bindings, nil
}
