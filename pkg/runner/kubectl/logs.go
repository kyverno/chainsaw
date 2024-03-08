package kubectl

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Logs(bindings binding.Bindings, collector *v1alpha1.PodLogs) (*v1alpha1.Command, error) {
	if collector == nil {
		return nil, errors.New("collector is null")
	}
	name, err := ConvertString(collector.Name, bindings)
	if err != nil {
		return nil, err
	}
	namespace, err := ConvertString(collector.Namespace, bindings)
	if err != nil {
		return nil, err
	}
	selector, err := ConvertString(collector.Selector, bindings)
	if err != nil {
		return nil, err
	}
	container, err := ConvertString(collector.Container, bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := ConvertString(collector.Cluster, bindings)
	if err != nil {
		return nil, err
	}
	if name == "" && selector == "" {
		return nil, errors.New("a name or selector must be specified")
	}
	if name != "" && selector != "" {
		return nil, errors.New("name cannot be provided when a selector is specified")
	}
	cmd := v1alpha1.Command{
		Cluster:    cluster,
		Timeout:    collector.Timeout,
		Entrypoint: "kubectl",
		Args:       []string{"logs", "--prefix"},
	}
	if name != "" {
		cmd.Args = append(cmd.Args, name)
	} else if selector != "" {
		cmd.Args = append(cmd.Args, "-l", selector)
	}
	if namespace == "" {
		namespace = "$NAMESPACE"
	}
	cmd.Args = append(cmd.Args, "-n", namespace)
	if container == "" {
		cmd.Args = append(cmd.Args, "--all-containers")
	} else {
		cmd.Args = append(cmd.Args, "-c", container)
	}
	if collector.Tail != nil {
		cmd.Args = append(cmd.Args, "--tail", fmt.Sprint(*collector.Tail))
	}
	return &cmd, nil
}
