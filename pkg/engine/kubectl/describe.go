package kubectl

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Describe(ctx context.Context, client client.Client, tc binding.Bindings, collector *v1alpha1.Describe) (string, []string, error) {
	if collector == nil {
		return "", nil, errors.New("collector is null")
	}
	name, err := collector.Name.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	namespace, err := collector.Namespace.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	selector, err := collector.Selector.Value(ctx, tc)
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
	args := []string{"describe", resource}
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
	if collector.ShowEvents != nil {
		args = append(args, fmt.Sprintf("--show-events=%t", *collector.ShowEvents))
	}
	return "kubectl", args, nil
}
