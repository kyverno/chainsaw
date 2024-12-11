package context

import (
	"context"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/expressions"
)

func WithBindings(tc TestContext, variables ...v1alpha1.Binding) (TestContext, error) {
	for _, variable := range variables {
		name, value, err := bindings.ResolveBinding(context.TODO(), tc.Compilers(), tc.Bindings(), nil, variable)
		if err != nil {
			return tc, err
		}
		tc = tc.WithBinding(name, value)
	}
	return tc, nil
}

func WithClusters(tc TestContext, basePath string, c map[string]v1alpha1.Cluster) TestContext {
	for name, cluster := range c {
		kubeconfig := filepath.Join(basePath, cluster.Kubeconfig)
		cluster := clusters.NewClusterFromKubeconfig(kubeconfig, cluster.Context)
		tc = tc.WithCluster(name, cluster)
	}
	return tc
}

func WithCurrentCluster(tc TestContext, name string) (TestContext, error) {
	name, err := expressions.String(context.TODO(), tc.Compilers(), name, tc.Bindings())
	if err != nil {
		return tc, err
	}
	tc = tc.WithCurrentCluster(name)
	config, client, err := tc.CurrentClusterClient()
	if err != nil {
		return tc, err
	}
	tc = tc.WithBinding("client", client)
	tc = tc.WithBinding("config", config)
	return tc, nil
}

func WithNamespace(tc TestContext, namespace string) TestContext {
	return tc.WithBinding("namespace", namespace)
}

func WithValues(tc TestContext, values any) TestContext {
	return tc.WithBinding("values", values)
}
