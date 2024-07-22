package model

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

type GlobalContext interface {
	Clusters() clusters.Registry
	Configuration() Configuration
	Namespace(context.Context) (*corev1.Namespace, error)
	TestContext(context.Context, *Test, int, int) TestContext
}

type TestContext interface {
	Bindings() binding.Bindings
	Clusters() clusters.Registry
	Configuration() Configuration
	// Namespace(context.Context) (*corev1.Namespace, error)
	// WithBindings(context.Context, string, any) GlobalContext
}

type globalContext struct {
	config   Configuration
	bindings binding.Bindings
	clusters clusters.Registry
}

type testContext struct {
	config   Configuration
	bindings binding.Bindings
	clusters clusters.Registry
}

func NewContext(ctx context.Context, values any, cluster *rest.Config, config Configuration) (GlobalContext, error) {
	tc := globalContext{
		config:   config,
		bindings: binding.NewBindings(),
		clusters: clusters.NewRegistry(),
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

func (tc *globalContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *globalContext) Configuration() Configuration {
	return tc.config
}

func (tc *globalContext) Namespace(ctx context.Context) (*corev1.Namespace, error) {
	if tc.config.Namespace.Name == "" {
		return nil, nil
	}
	return buildNamespace(ctx, tc.config.Namespace.Name, tc.config.Namespace.Template, tc.bindings)
}

func (tc *globalContext) TestContext(ctx context.Context, test *Test, i int, s int) TestContext {
	return &testContext{
		config:   tc.config,
		clusters: tc.clusters,
		bindings: apibindings.RegisterNamedBinding(ctx, tc.bindings, "test",
			TestInfo{
				Id:         i + 1,
				ScenarioId: s + 1,
				Metadata:   test.ObjectMeta,
			},
		),
	}
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
