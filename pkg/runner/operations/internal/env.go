package internal

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
)

func RegisterEnvs(ctx context.Context, namespace string, bindings binding.Bindings, envs ...v1alpha1.Binding) (map[string]string, []string, error) {
	mapOut := map[string]string{}
	var envsOut []string
	for _, env := range envs {
		if err := env.CheckEnvName(); err != nil {
			return mapOut, envsOut, err
		}
		patched, err := mutation.Mutate(ctx, nil, mutation.Parse(ctx, env.Value.Value), nil, bindings, template.WithFunctionCaller(functions.Caller))
		if err != nil {
			return mapOut, envsOut, err
		}
		if patched, ok := patched.(string); !ok {
			return mapOut, envsOut, fmt.Errorf("value must be a string (%s)", env.Name)
		} else {
			mapOut[env.Name] = patched
			envsOut = append(envsOut, env.Name+"="+patched)
		}
	}
	mapOut["NAMESPACE"] = namespace
	envsOut = append(envsOut, "NAMESPACE="+namespace)
	return mapOut, envsOut, nil
}
