package internal

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	apibindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

func RegisterEnvs(ctx context.Context, compilers compilers.Compilers, namespace string, bindings apis.Bindings, envs ...v1alpha1.Binding) (map[string]string, []string, error) {
	mapOut := map[string]string{}
	var envsOut []string
	for _, env := range envs {
		name, value, err := apibindings.ResolveBinding(ctx, compilers, bindings, nil, env)
		if err != nil {
			return mapOut, envsOut, err
		}
		if value, ok := value.(string); !ok {
			return mapOut, envsOut, fmt.Errorf("value must be a string (%s)", env.Name)
		} else {
			mapOut[name] = value
			envsOut = append(envsOut, name+"="+value)
		}
	}
	mapOut["NAMESPACE"] = namespace
	envsOut = append(envsOut, "NAMESPACE="+namespace)
	return mapOut, envsOut, nil
}
