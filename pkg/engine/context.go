package engine

import (
	"context"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
)

func WithBindings(ctx context.Context, tc Context, variables ...v1alpha1.Binding) (Context, error) {
	for _, variable := range variables {
		name, value, err := bindings.ResolveBinding(ctx, tc.Bindings(), nil, variable)
		if err != nil {
			return tc, err
		}
		tc = tc.WithBinding(ctx, name, value)
	}
	return tc, nil
}

func WithClusters(ctx context.Context, tc Context, basePath string, c map[string]v1alpha1.Cluster) Context {
	for name, cluster := range c {
		kubeconfig := filepath.Join(basePath, cluster.Kubeconfig)
		cluster := clusters.NewClusterFromKubeconfig(kubeconfig, cluster.Context)
		tc = tc.WithCluster(ctx, name, cluster)
	}
	return tc
}

func WithCurrentCluster(ctx context.Context, tc Context, name string) (Context, error) {
	tc = tc.WithCurrentCluster(ctx, name)
	config, client, err := tc.CurrentClusterClient()
	if err != nil {
		return tc, err
	}
	tc = tc.WithBinding(ctx, "client", client)
	tc = tc.WithBinding(ctx, "config", config)
	return tc, nil
}

func WithNamespace(ctx context.Context, tc Context, namespace string) Context {
	return tc.WithBinding(ctx, "namespace", namespace)
}

func WithValues(ctx context.Context, tc Context, values any) Context {
	return tc.WithBinding(ctx, "values", values)
}
