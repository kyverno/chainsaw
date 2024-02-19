package processors

import (
	"github.com/kyverno/chainsaw/pkg/client"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"k8s.io/client-go/rest"
)

const DefaultClient = ""

type clusters struct {
	clients map[string]client.Client
}

func NewClusters() clusters {
	return clusters{
		clients: map[string]client.Client{},
	}
}

func (c *clusters) Register(name string, cfg *rest.Config) (client.Client, error) {
	client, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	c.clients[DefaultClient] = runnerclient.New(client)
	return client, nil
}

func (c *clusters) client(names ...string) client.Client {
	for _, name := range names {
		if name != "" {
			return c.clients[name]
		}
	}
	return c.clients[DefaultClient]
}
