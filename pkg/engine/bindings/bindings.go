package bindings

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

var identifier = regexp.MustCompile(`^\w+$`)

func checkBindingName(name string) error {
	if !identifier.MatchString(name) {
		return fmt.Errorf("invalid binding name %s", name)
	}
	return nil
}

func RegisterBinding(ctx context.Context, bindings binding.Bindings, name string, value any) binding.Bindings {
	return bindings.Register("$"+name, binding.NewBinding(value))
}

func ResolveBinding(ctx context.Context, bindings binding.Bindings, input any, variable v1alpha1.Binding) (string, any, error) {
	name, err := String(variable.Name, bindings)
	if err != nil {
		return "", nil, err
	}
	if err := checkBindingName(name); err != nil {
		return "", nil, err
	}
	value, err := mutation.Mutate(ctx, nil, mutation.Parse(ctx, variable.Value.Value), input, bindings, template.WithFunctionCaller(functions.Caller))
	if err != nil {
		return "", nil, err
	}
	return name, value, err
}
