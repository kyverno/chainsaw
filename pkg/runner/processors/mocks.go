package processors

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"k8s.io/client-go/rest"
)

type registryMock struct {
	client client.Client
}

func (r registryMock) Register(string, clusters.Cluster) clusters.Registry {
	return r
}

func (r registryMock) Lookup(string) clusters.Cluster {
	return nil
}

func (r registryMock) Build(clusters.Cluster) (*rest.Config, client.Client, error) {
	return nil, r.client, nil
}
