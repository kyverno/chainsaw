package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
)

func Client(cluster Cluster) (client.Client, error) {
	client, err := client.New(cluster.Config())
	if err != nil {
		return nil, err
	}
	return client, nil
}
