package runner

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func collect(collector v1alpha1.Collect) []v1alpha1.Command {
	var commands []v1alpha1.Command
	if collector.Events != nil {
		commands = append(commands, events(*collector.Events))
	}
	if collector.PodLogs != nil {
		commands = append(commands, podLogs(*collector.PodLogs))
	}
	return commands
}

func events(collector v1alpha1.Events) v1alpha1.Command {
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
	return cmd
}

func podLogs(collector v1alpha1.PodLogs) v1alpha1.Command {
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
	return cmd
}
