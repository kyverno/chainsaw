package model

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"k8s.io/client-go/rest"
)

type TestContext interface {
	Bindings() binding.Bindings
	Clusters() clusters.Registry
	Configuration() Configuration
	Cluster() (*rest.Config, client.Client, error)
	WithBindings(binding.Bindings) TestContext
}

type testContext struct {
	config   Configuration
	bindings binding.Bindings
	clusters clusters.Registry
	cluster  string
}

func NewContext(ctx context.Context, values any, cluster *rest.Config, config Configuration) (TestContext, error) {
	tc := testContext{
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
			return nil, err
		}
		tc.clusters = tc.clusters.Register(clusters.DefaultClient, cluster)
		// register default cluster in bindings
		clusterConfig, clusterClient, err := tc.clusters.Resolve(false)
		if err != nil {
			return nil, err
		}
		tc.bindings = apibindings.RegisterClusterBindings(ctx, tc.bindings, clusterConfig, clusterClient)
	}
	// 3. register clusters
	tc.clusters = clusters.Register(tc.clusters, "", config.Clusters)
	return &tc, nil
}

func (tc *testContext) Bindings() binding.Bindings {
	return tc.bindings
}

func (tc *testContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *testContext) Cluster() (*rest.Config, client.Client, error) {
	return tc.clusters.Resolve(false, tc.cluster)
}

func (tc *testContext) Configuration() Configuration {
	return tc.config
}

func (tc *testContext) WithBindings(bindings binding.Bindings) TestContext {
	return &testContext{
		config:   tc.config,
		bindings: bindings,
		clusters: tc.clusters,
		cluster:  tc.cluster,
	}
}
