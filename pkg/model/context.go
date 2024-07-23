package model

import (
	"context"
	"path/filepath"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"k8s.io/client-go/rest"
)

type TestContext struct {
	bindings binding.Bindings
	clusters clusters.Registry
	cluster  string
}

func MakeContext(bindings binding.Bindings, registry clusters.Registry) TestContext {
	return TestContext{
		bindings: bindings,
		clusters: registry,
		cluster:  clusters.DefaultClient,
	}
}

func EmptyContext() TestContext {
	return MakeContext(binding.NewBindings(), clusters.NewRegistry())
}

func (tc *TestContext) Bindings() binding.Bindings {
	return tc.bindings
}

func (tc *TestContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *TestContext) Cluster() (*rest.Config, client.Client, error) {
	return tc.clusters.Resolve(false, tc.cluster)
}

func (tc TestContext) WithCluster(ctx context.Context, name string, cluster clusters.Cluster) TestContext {
	tc.clusters = tc.clusters.Register(name, cluster)
	return tc
}

func (tc TestContext) WithBinding(ctx context.Context, name string, value any) TestContext {
	tc.bindings = apibindings.RegisterNamedBinding(ctx, tc.bindings, name, value)
	return tc
}

func WithValues(ctx context.Context, tc TestContext, values any) TestContext {
	return tc.WithBinding(ctx, "values", values)
}

func WithClusters(ctx context.Context, tc TestContext, basePath string, c map[string]v1alpha1.Cluster) TestContext {
	for name, cluster := range c {
		kubeconfig := filepath.Join(basePath, cluster.Kubeconfig)
		cluster := clusters.NewClusterFromKubeconfig(kubeconfig, cluster.Context)
		tc = tc.WithCluster(ctx, name, cluster)
	}
	return tc
}

func UseCluster(ctx context.Context, tc TestContext, name string) (TestContext, error) {
	clusterConfig, clusterClient, err := tc.clusters.Resolve(false)
	if err != nil {
		return tc, err
	}
	tc = tc.WithBinding(ctx, "client", clusterClient)
	tc = tc.WithBinding(ctx, "config", clusterConfig)
	return tc, nil
}
