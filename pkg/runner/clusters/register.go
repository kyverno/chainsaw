package clusters

import (
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Register(registry Registry, basePath string, clusters map[string]v1alpha1.Cluster) Registry {
	for name, cluster := range clusters {
		registry = registry.Register(name, NewClusterFromKubeconfig(filepath.Join(basePath, cluster.Kubeconfig), cluster.Context))
	}
	return registry
}
