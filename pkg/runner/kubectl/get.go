package kubectl

import (
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Get(client client.Client, bindings binding.Bindings, collector *v1alpha1.Get) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(client, bindings, collector.ResourceReference)
	if err != nil {
		return nil, err
	}
	cmd := v1alpha1.Command{
		Cluster:    collector.Cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"get", resource},
	}
	if collector.Name != "" {
		cmd.Args = append(cmd.Args, collector.Name)
	} else if collector.Selector != "" {
		cmd.Args = append(cmd.Args, "-l", collector.Selector)
	}
	if !clustered {
		if collector.Namespace == "*" {
			cmd.Args = append(cmd.Args, "--all-namespaces")
		} else {
			namespace := collector.Namespace
			if namespace == "" {
				namespace = "$NAMESPACE"
			}
			cmd.Args = append(cmd.Args, "-n", namespace)
		}
	}
	if collector.Format != "" {
		cmd.Args = append(cmd.Args, "-o", string(collector.Format))
	}
	return &cmd, nil
}
