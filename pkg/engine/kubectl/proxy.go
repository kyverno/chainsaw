package kubectl

import (
	"errors"
	"path"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
)

func Proxy(client client.Client, bindings binding.Bindings, collector *v1alpha1.Proxy) (*v1alpha1.Command, error) {
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
	targetPath, err := apibindings.String(collector.TargetPath, bindings)
	if err != nil {
		return nil, err
	}
	targetPort, err := apibindings.String(collector.TargetPort, bindings)
	if err != nil {
		return nil, err
	}
	cluster, err := apibindings.String(collector.Cluster, bindings)
	if err != nil {
		return nil, err
	}
	resource, _, err := mapResource(client, bindings, collector.ObjectType)
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
		ActionOutputs: collector.ActionOutputs,
		Entrypoint:    "kubectl",
		Args:          []string{"get", "--raw", path.Join("/api", "v1", "namespaces", namespace, resource, name+":"+targetPort, "proxy", targetPath)},
	}
	return &cmd, nil
}
