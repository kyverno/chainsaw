package model

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

type TestContext interface {
	Bindings() binding.Bindings
	Clusters() clusters.Registry
	Configuration() Configuration
	Namespace(context.Context) (*corev1.Namespace, error)
	WithBindings(ctx context.Context, name string, value any) TestContext
}

type testContext struct {
	config   Configuration
	bindings binding.Bindings
	clusters clusters.Registry
}

func NewContext(ctx context.Context, values any, cluster *rest.Config, config Configuration) (TestContext, error) {
	out := testContext{
		config:   config,
		bindings: binding.NewBindings(),
		clusters: clusters.NewRegistry(),
	}
	// 1. register values first
	out.bindings = apibindings.RegisterNamedBinding(ctx, out.bindings, "values", values)
	// 2. register default cluster
	if cluster != nil {
		cluster, err := clusters.NewClusterFromConfig(cluster)
		if err != nil {
			return nil, err
		}
		out.clusters = out.clusters.Register(clusters.DefaultClient, cluster)
		// register default cluster in bindings
		clusterConfig, clusterClient, err := out.clusters.Resolve(false)
		if err != nil {
			return nil, err
		}
		out.bindings = apibindings.RegisterClusterBindings(ctx, out.bindings, clusterConfig, clusterClient)
	}
	// 3. register clusters
	out.clusters = clusters.Register(out.clusters, "", config.Clusters)
	return &out, nil
}

func (tc *testContext) Bindings() binding.Bindings {
	return tc.bindings
}

func (tc *testContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *testContext) Configuration() Configuration {
	return tc.config
}

func (tc *testContext) Namespace(ctx context.Context) (*corev1.Namespace, error) {
	if tc.config.Namespace.Name == "" {
		return nil, nil
	}
	return buildNamespace(ctx, tc.config.Namespace.Name, tc.config.Namespace.Template, tc.bindings)
}

func (tc *testContext) WithBindings(ctx context.Context, name string, value any) TestContext {
	return &testContext{
		config:   tc.config,
		clusters: tc.clusters,
		bindings: apibindings.RegisterNamedBinding(ctx, tc.bindings, name, value),
	}
}
