package kubectl

import (
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
)

func Get(client client.Client, bindings binding.Bindings, collector *v1alpha1.Get) (*v1alpha1.Command, error) {
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
	format, err := apibindings.String(string(collector.Format), bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := apibindings.StringPointer(collector.Cluster, bindings)
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
		Args:       []string{"get", resource},
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
	if format != "" {
		cmd.Args = append(cmd.Args, "-o", format)
	}
	return &cmd, nil
}
