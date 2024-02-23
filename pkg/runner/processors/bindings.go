package processors

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	identifier     = regexp.MustCompile(`^\w+$`)
	forbiddenNames = []string{"namespace", "client", "error", "values", "stdout", "stderr"}
	forbidden      = sets.New(forbiddenNames...)
)

func registerBindings(ctx context.Context, bindings binding.Bindings, variables ...v1alpha1.Binding) (binding.Bindings, error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	for _, variable := range variables {
		if err := checkBindingName(variable.Name); err != nil {
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

func checkBindingName(name string) error {
	if forbidden.Has(name) {
		return fmt.Errorf("binding name is forbidden (%s), it must not be (%s)", name, strings.Join(forbiddenNames, ", "))
	}
	if !identifier.MatchString(name) {
		return fmt.Errorf("invalid binding name %s", name)
	}
	return nil
}
