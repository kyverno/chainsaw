package kubectl

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Describe(client client.Client, bindings binding.Bindings, collector *v1alpha1.Describe) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	name, err := convertString(collector.Name, bindings)
	if err != nil {
		return nil, err
	}
	namespace, err := convertString(collector.Namespace, bindings)
	if err != nil {
		return nil, err
	}
	selector, err := convertString(collector.Selector, bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := convertString(collector.Cluster, bindings)
	if err != nil {
		return nil, err
	}
	if name != "" && selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(client, bindings, collector.ResourceReference)
	if err != nil {
		return nil, err
	}
	cmd := v1alpha1.Command{
		Cluster:    cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"describe", resource},
	}
	if name != "" {
		cmd.Args = append(cmd.Args, name)
	} else if selector != "" {
		cmd.Args = append(cmd.Args, "-l", selector)
	}
	if !clustered {
		if namespace == "*" {
			cmd.Args = append(cmd.Args, "--all-namespaces")
		} else {
			if namespace == "" {
				namespace = "$NAMESPACE"
			}
			cmd.Args = append(cmd.Args, "-n", namespace)
		}
	}
	if collector.ShowEvents != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--show-events=%t", *collector.ShowEvents))
	}
	return &cmd, nil
}
