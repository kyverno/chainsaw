package kubectl

import (
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Get(collector *v1alpha1.Get) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, nil
	}
	if collector.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	cmd := v1alpha1.Command{
		Entrypoint: "kubectl",
		Args:       []string{"get", collector.Resource},
	}
	if collector.Name != "" {
		cmd.Args = append(cmd.Args, collector.Name)
	}
	if collector.Selector != "" {
		cmd.Args = append(cmd.Args, "-l", collector.Selector)
	}
	// TODO: what if cluster scoped resource ?
	namespace := collector.Namespace
	if collector.Namespace == "" {
		namespace = "$NAMESPACE"
	}
	if collector.Format != "" {
		cmd.Args = append(cmd.Args, "-o", string(collector.Format))
	}
	cmd.Args = append(cmd.Args, "-n", namespace)
	return &cmd, nil
}
