package clusters

import (
	"sync"

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
	resolver func() (*rest.Config, error)
}

func NewClusterFromKubeconfig(kubeconfig string, context string) Cluster {
	resolver := sync.OnceValues(func() (*rest.Config, error) {
		return restutils.Config(kubeconfig, clientcmd.ConfigOverrides{
			CurrentContext: context,
		})
	})
	return &fromKubeconfig{
		resolver: resolver,
	}
}

func (c *fromKubeconfig) Config() (*rest.Config, error) {
	return c.resolver()
}
