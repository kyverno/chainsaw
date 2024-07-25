package kubectl

import (
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
)

func Logs(bindings binding.Bindings, collector *v1alpha1.PodLogs) (*v1alpha1.Command, error) {
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
	container, err := apibindings.String(collector.Container, bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := apibindings.StringPointer(collector.Cluster, bindings)
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
		ActionClusters: v1alpha1.ActionClusters{
			Cluster:  cluster,
			Clusters: collector.Clusters,
		},
		ActionTimeout: v1alpha1.ActionTimeout{
			Timeout: collector.Timeout,
		},
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
