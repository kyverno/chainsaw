package kubectl

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
)

func Describe(client client.Client, bindings binding.Bindings, collector *v1alpha1.Describe) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	name, err := apibindings.String(collector.Name, bindings)
	if err != nil {
		return nil, err
	}
	namespace, err := apibindings.String(collector.Namespace, bindings)
	if err != nil {
		return nil, err
	}
	selector, err := apibindings.String(collector.Selector, bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := apibindings.String(collector.Cluster, bindings)
	if err != nil {
		return nil, err
	}
	if name != "" && selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	resource, clustered, err := mapResource(client, bindings, collector.ObjectType)
	if err != nil {
		return nil, err
	}
	cmd := v1alpha1.Command{
		ActionClusters: v1alpha1.ActionClusters{
			Cluster:  cluster,
			Clusters: collector.Clusters,
		},
		ActionTimeout: v1alpha1.ActionTimeout{
			Timeout: collector.Timeout,
		},
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
