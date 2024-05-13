package clusters

import (
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Cluster interface {
	Config() (*rest.Config, error)
}

type fromConfig struct {
	config *rest.Config
}

func NewClusterFromConfig(config *rest.Config) (Cluster, error) {
	return &fromConfig{
		config: config,
	}, nil
}

func (c *fromConfig) Config() (*rest.Config, error) {
	return c.config, nil
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

func (c *fromKubeconfig) Config() (*rest.Config, error) {
	return restutils.Config(c.kubeconfig, clientcmd.ConfigOverrides{
		CurrentContext: c.context,
	})
}
