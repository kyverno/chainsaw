package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/client-go/rest"
)

func link(cluster Cluster, dryRun bool) (*rest.Config, client.Client, error) {
	config, err := cluster.Config()
	if err != nil {
		return nil, nil, err
	}
	client, err := makeClient(config, dryRun)
	if err != nil {
		return nil, nil, err
	}
	return config, client, nil
}
