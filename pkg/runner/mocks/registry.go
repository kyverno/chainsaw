package mocks

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"k8s.io/client-go/rest"
)

type Registry struct {
	Client client.Client
}

func (r Registry) Register(string, clusters.Cluster) clusters.Registry {
	return r
}

func (r Registry) Lookup(string) clusters.Cluster {
	return nil
}

func (r Registry) Build(clusters.Cluster) (*rest.Config, client.Client, error) {
	return nil, r.Client, nil
}
