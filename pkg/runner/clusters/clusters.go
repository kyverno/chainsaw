package processors

import (
	"github.com/kyverno/chainsaw/pkg/client"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"k8s.io/client-go/rest"
)

const DefaultClient = ""

type cluster struct {
	config *rest.Config
	client client.Client
}

type clusters struct {
	clients map[string]cluster
}

func New() clusters {
	return clusters{
		clients: map[string]cluster{},
	}
}

func (c *clusters) Register(name string, config *rest.Config) error {
	client, err := client.New(config)
	if err != nil {
		return err
	}
	c.clients[DefaultClient] = cluster{
		config: config,
		client: runnerclient.New(client),
	}
	return nil
}

func (c *clusters) client(names ...string) (*rest.Config, client.Client) {
	for _, name := range names {
		if name != "" {
			cluster := c.clients[name]
			return cluster.config, cluster.client
		}
	}
	cluster := c.clients[DefaultClient]
	return cluster.config, cluster.client
}
