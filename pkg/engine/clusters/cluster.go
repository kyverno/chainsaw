package clusters

import (
	"sync"

	"github.com/kyverno/chainsaw/pkg/model"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Cluster interface {
	Config() (*rest.Config, error)
	model.WarningsHolder
}

type fromConfig struct {
	config *rest.Config
	*model.WithWarnings
}

func NewClusterFromConfig(config *rest.Config) Cluster {
	return &fromConfig{
		config:       config,
		WithWarnings: model.NewWithWarnings(),
	}
}

func (c *fromConfig) Config() (*rest.Config, error) {
	return c.config, nil
}

type fromKubeconfig struct {
	resolver func() (*rest.Config, error)
	*model.WithWarnings
}

func NewClusterFromKubeconfig(kubeconfig string, context string) Cluster {
	resolver := sync.OnceValues(func() (*rest.Config, error) {
		return restutils.Config(kubeconfig, clientcmd.ConfigOverrides{
			CurrentContext: context,
		})
	})
	return &fromKubeconfig{
		resolver:     resolver,
		WithWarnings: model.NewWithWarnings(),
	}
}

func (c *fromKubeconfig) Config() (*rest.Config, error) {
	return c.resolver()
}
