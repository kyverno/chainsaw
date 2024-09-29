package bindings

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

var identifier = regexp.MustCompile(`^\w+$`)

func checkBindingName(name string) error {
	if !identifier.MatchString(name) {
		return fmt.Errorf("invalid binding name %s", name)
	}
	return nil
}

func RegisterBinding(ctx context.Context, bindings apis.Bindings, name string, value any) apis.Bindings {
	return bindings.Register("$"+name, apis.NewBinding(value))
}

func ResolveBinding(ctx context.Context, compilers compilers.Compilers, bindings apis.Bindings, input any, variable v1alpha1.Binding) (string, any, error) {
	name, err := variable.Name.Value(ctx, compilers, bindings)
	if err != nil {
		return "", nil, err
	}
	if err := checkBindingName(name); err != nil {
		return "", nil, err
	}
	if variable.Compiler != nil {
		compilers = compilers.WithDefaultCompiler(string(*variable.Compiler))
	}
	value, err := templating.Template(ctx, compilers, variable.Value, input, bindings)
	if err != nil {
		return "", nil, err
	}
	return name, value, err
}
