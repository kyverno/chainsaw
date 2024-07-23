package model

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"k8s.io/client-go/rest"
)

type TestContext struct {
	config   Configuration
	bindings binding.Bindings
	clusters clusters.Registry
	cluster  string
}

func EmptyContext(config Configuration) TestContext {
	return TestContext{
		config:   config,
		bindings: binding.NewBindings(),
		clusters: clusters.NewRegistry(),
		cluster:  clusters.DefaultClient,
	}
}

func MakeContext(config Configuration, bindings binding.Bindings, registry clusters.Registry) TestContext {
	return TestContext{
		config:   config,
		bindings: bindings,
		clusters: registry,
		cluster:  clusters.DefaultClient,
	}
}

func NewContext(ctx context.Context, values any, cluster *rest.Config, config Configuration) (TestContext, error) {
	tc := TestContext{
		config:   config,
		bindings: binding.NewBindings(),
		clusters: clusters.NewRegistry(),
		cluster:  clusters.DefaultClient,
	}
	// 1. register values first
	tc.bindings = apibindings.RegisterNamedBinding(ctx, tc.bindings, "values", values)
	// 2. register default cluster
	if cluster != nil {
		cluster, err := clusters.NewClusterFromConfig(cluster)
		if err != nil {
			return tc, err
		}
		tc.clusters = tc.clusters.Register(clusters.DefaultClient, cluster)
		// register default cluster in bindings
		clusterConfig, clusterClient, err := tc.clusters.Resolve(false)
		if err != nil {
			return tc, err
		}
		tc.bindings = apibindings.RegisterClusterBindings(ctx, tc.bindings, clusterConfig, clusterClient)
	}
	// 3. register clusters
	tc.clusters = clusters.Register(tc.clusters, "", config.Clusters)
	return tc, nil
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

func (tc *TestContext) Configuration() Configuration {
	return tc.config
}

func (tc TestContext) WithBindings(bindings binding.Bindings) TestContext {
	return TestContext{
		config:   tc.config,
		bindings: bindings,
		clusters: tc.clusters,
		cluster:  tc.cluster,
	}
}

func (tc TestContext) WithBinding(ctx context.Context, name string, value any) TestContext {
	tc.bindings = apibindings.RegisterNamedBinding(ctx, tc.bindings, name, value)
	return tc
}

func (tc TestContext) WithValues(ctx context.Context, values any) TestContext {
	return tc.WithBinding(ctx, "values", values)
}
