package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/client-go/rest"
)

const DefaultClient = ""

type Registry interface {
	Register(string, Cluster) Registry
	Resolve(bool, ...string) (*rest.Config, client.Client, error)
}

type registry struct {
	clusters map[string]Cluster
}

func NewRegistry() Registry {
	return registry{
		clusters: map[string]Cluster{},
	}
}

func (c registry) Register(name string, cluster Cluster) Registry {
	values := map[string]Cluster{}
	for k, v := range c.clusters {
		values[k] = v
	}
	values[name] = cluster
	return registry{
		clusters: values,
	}
}

func (c registry) find(names ...string) Cluster {
	for _, name := range names {
		if name != "" {
			return c.clusters[name]
		}
	}
	return c.clusters[DefaultClient]
}

func (c registry) Resolve(dryRun bool, names ...string) (*rest.Config, client.Client, error) {
	cluster := c.find(names...)
	if cluster != nil {
		return link(cluster, dryRun)
	}
	return nil, nil, nil
}
