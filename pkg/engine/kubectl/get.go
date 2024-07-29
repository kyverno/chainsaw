package kubectl

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
)

func Get(ctx context.Context, client client.Client, tc binding.Bindings, collector *v1alpha1.Get) (string, []string, error) {
	if collector == nil {
		return "", nil, errors.New("collector is null")
	}
	name, err := templating.String(ctx, collector.Name, tc)
	if err != nil {
		return "", nil, err
	}
	namespace, err := templating.String(ctx, collector.Namespace, tc)
	if err != nil {
		return "", nil, err
	}
	selector, err := templating.String(ctx, collector.Selector, tc)
	if err != nil {
		return "", nil, err
	}
	format, err := templating.String(ctx, string(collector.Format), tc)
	if err != nil {
		return "", nil, err
	}
	if name != "" && selector != "" {
		return "", nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(ctx, client, tc, collector.ObjectType)
	if err != nil {
		return "", nil, err
	}
	args := []string{"get", resource}
	if name != "" {
		args = append(args, name)
	} else if selector != "" {
		args = append(args, "-l", selector)
	}
	if !clustered {
		if namespace == "*" {
			args = append(args, "--all-namespaces")
		} else {
			if namespace == "" {
				namespace = "$NAMESPACE"
			}
			args = append(args, "-n", namespace)
		}
	}
	if format != "" {
		args = append(args, "-o", format)
	}
	return "kubectl", args, nil
}
