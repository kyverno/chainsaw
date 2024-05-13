package clusters

const DefaultClient = ""

type Registry interface {
	Register(string, Cluster) Registry
	Resolve(...string) Cluster
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

func (c registry) Resolve(names ...string) Cluster {
	for _, name := range names {
		if name != "" {
			return c.clusters[name]
		}
	}
	return c.clusters[DefaultClient]
}
