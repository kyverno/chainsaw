package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	runnerclient "github.com/kyverno/chainsaw/pkg/runner/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Cluster interface {
	Config() *rest.Config
	Client() client.Client
}

type fromConfig struct {
	config *rest.Config
	client client.Client
}

func NewClusterFromConfig(config *rest.Config) (Cluster, error) {
	client, err := client.New(config)
	if err != nil {
		return nil, err
	}
	return &fromConfig{
		config: config,
		client: runnerclient.New(client),
	}, nil
}

func (c *fromConfig) Config() *rest.Config {
	return c.config
}

func (c *fromConfig) Client() client.Client {
	return c.client
}

type fromKubeconfig struct {
	kubeconfig string
	context    string
}

func NewClusterFromKubeconfig(kubeconfig string, context string) Cluster {
	return &fromKubeconfig{
		kubeconfig: kubeconfig,
		context:    context,
	}
}

func (c *fromKubeconfig) Config() *rest.Config {
	cfg, err := restutils.Config(c.kubeconfig, clientcmd.ConfigOverrides{
		CurrentContext: c.context,
	})
	if err != nil {
		// TODO: error handling
		panic(err)
	}
	return cfg
}

func (c *fromKubeconfig) Client() client.Client {
	client, err := client.New(c.Config())
	if err != nil {
		// TODO: error handling
		panic(err)
	}
	return runnerclient.New(client)
}

type dryRun struct {
	inner Cluster
}

func NewDryRun(inner Cluster) Cluster {
	return &dryRun{
		inner: inner,
	}
}

func (c *dryRun) Config() *rest.Config {
	return c.inner.Config()
}

func (c *dryRun) Client() client.Client {
	return client.DryRun(c.inner.Client())
}
