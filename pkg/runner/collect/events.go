package collect

import (
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func events(collector *v1alpha1.Events) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, nil
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	cmd := v1alpha1.Command{
		Entrypoint: "kubectl",
		Args:       []string{"get", "events"},
	}
	if collector.Name != "" {
		cmd.Args = append(cmd.Args, collector.Name)
	}
	if collector.Selector != "" {
		cmd.Args = append(cmd.Args, "-l", collector.Selector)
	}
	namespace := collector.Namespace
	if collector.Namespace == "" {
		namespace = "$NAMESPACE"
	}
	cmd.Args = append(cmd.Args, "-n", namespace)
	return &cmd, nil
}
