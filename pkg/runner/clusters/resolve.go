package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/client-go/rest"
)

func Link(cluster Cluster, dryRun bool) (*rest.Config, client.Client, error) {
	config, err := cluster.Config()
	if err != nil {
		return nil, nil, err
	}
	client, err := Client(config, dryRun)
	if err != nil {
		return nil, nil, err
	}
	return config, client, nil
}

func Resolve(registry Registry, names ...string) (*rest.Config, client.Client, error) {
	cluster := registry.Resolve(names...)
	if cluster != nil {
		return Link(cluster, false)
	}
	return nil, nil, nil
}
