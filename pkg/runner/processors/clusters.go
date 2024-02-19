package processors

import (
	"github.com/kyverno/chainsaw/pkg/client"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	"k8s.io/client-go/rest"
)

const defaultClient = ""

type clusters struct {
	clients map[string]client.Client
}

func NewClusters() clusters {
	return clusters{
		clients: map[string]client.Client{},
	}
}

func (c *clusters) Register(cfg *rest.Config) (client.Client, error) {
	client, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	c.clients[defaultClient] = runnerclient.New(client)
	return client, nil
}

func (c *clusters) client(name string) client.Client {
	if name == "" {
		name = defaultClient
	}
	return c.clients[name]
}
