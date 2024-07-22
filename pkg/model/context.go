package model

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

type TestContext struct {
	config   Configuration
	bindings binding.Bindings
	clusters clusters.Registry
}

func NewContext(ctx context.Context, values any, cluster *rest.Config, config Configuration) (*TestContext, error) {
	out := TestContext{
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

func (tc *TestContext) Bindings() binding.Bindings {
	return tc.bindings
}

func (tc *TestContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *TestContext) Configuration() Configuration {
	return tc.config
}

func (tc *TestContext) Namespace(ctx context.Context) (*corev1.Namespace, error) {
	if tc.config.Namespace.Name == "" {
		return nil, nil
	}
	return buildNamespace(ctx, tc.config.Namespace.Name, tc.config.Namespace.Template, tc.bindings)
}
