package collect

import (
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func podLogs(collector *v1alpha1.PodLogs) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, nil
	}
	if collector.Name == "" && collector.Selector == "" {
		return nil, errors.New("a name or selector must be specified")
	}
	if collector.Name != "" && collector.Selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	cmd := v1alpha1.Command{
		Entrypoint: "kubectl",
		Args:       []string{"logs", "--prefix"},
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
	if collector.Container == "" {
		cmd.Args = append(cmd.Args, "--all-containers")
	} else {
		cmd.Args = append(cmd.Args, "-c", collector.Container)
	}
	if collector.Tail != nil {
		cmd.Args = append(cmd.Args, "--tail", fmt.Sprint(*collector.Tail))
	}
	return &cmd, nil
}
