package kubectl

import (
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Describe(collector *v1alpha1.Describe) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, nil
	}
	if collector.Resource == "" {
		return nil, errors.New("a resource must be specified")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	if collector.Name == "" && collector.Selector == "" {
		return nil, errors.New("a name or label selector must be specified")
	}
	cmd := v1alpha1.Command{
		Cluster:    collector.Cluster,
		Entrypoint: "kubectl",
		Args:       []string{"describe", collector.Resource},
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
	cmd.Args = append(cmd.Args, "-n", namespace)
	if collector.ShowEvents != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--show-events=%t", *collector.ShowEvents))
	}
	if collector.Timeout != nil {
		cmd.Timeout = collector.Timeout
	}
	return &cmd, nil
}
