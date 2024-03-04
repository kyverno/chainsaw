package kubectl

import (
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Describe(client client.Client, collector *v1alpha1.Describe) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(client, collector.ResourceReference)
	if err != nil {
		return nil, err
	}
	cmd := v1alpha1.Command{
		Cluster:    collector.Cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"describe", resource},
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
	if collector.ShowEvents != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--show-events=%t", *collector.ShowEvents))
	}
	return &cmd, nil
}
