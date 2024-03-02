package kubectl

import (
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/api/meta"
)

func Get(client client.Client, collector *v1alpha1.Get) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	if collector.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	cmd := v1alpha1.Command{
		Cluster:    collector.Cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"get", collector.Resource},
	}
	if collector.Name != "" {
		cmd.Args = append(cmd.Args, collector.Name)
	}
	if collector.Selector != "" {
		cmd.Args = append(cmd.Args, "-l", collector.Selector)
	}
	mapping, err := getMapping(client, collector.Resource)
	if err != nil {
		return nil, err
	}
	clustered := mapping.Scope.Name() == meta.RESTScopeNameRoot
	if !clustered {
		namespace := collector.Namespace
		if collector.Namespace == "" {
			namespace = "$NAMESPACE"
		}
		cmd.Args = append(cmd.Args, "-n", namespace)
	}
	if collector.Format != "" {
		cmd.Args = append(cmd.Args, "-o", string(collector.Format))
	}
	return &cmd, nil
}
