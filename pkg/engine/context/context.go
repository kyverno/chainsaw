package context

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/client/dryrun"
	apibindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/model"
	"k8s.io/client-go/rest"
)

type TestContext struct {
	*model.Summary
	*model.Report
	bindings binding.Bindings
	cluster  clusters.Cluster
	clusters clusters.Registry
	dryRun   bool
}

func MakeContext(bindings binding.Bindings, registry clusters.Registry) TestContext {
	return TestContext{
		Summary:  &model.Summary{},
		Report:   &model.Report{},
		bindings: bindings,
		clusters: registry,
		cluster:  nil,
	}
}

func EmptyContext() TestContext {
	return MakeContext(binding.NewBindings(), clusters.NewRegistry(nil))
}

func (tc *TestContext) Bindings() binding.Bindings {
	return tc.bindings
}

func (tc *TestContext) Cluster(name string) clusters.Cluster {
	return tc.clusters.Lookup(name)
}

func (tc *TestContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *TestContext) CurrentCluster() clusters.Cluster {
	return tc.cluster
}

func (tc *TestContext) CurrentClusterClient() (*rest.Config, client.Client, error) {
	config, client, err := tc.clusters.Build(tc.cluster)
	if err == nil && client != nil && tc.DryRun() {
		client = dryrun.New(client)
	}
	return config, client, err
}

func (tc *TestContext) DryRun() bool {
	return tc.dryRun
}

func (tc TestContext) WithBinding(ctx context.Context, name string, value any) TestContext {
	tc.bindings = apibindings.RegisterBinding(ctx, tc.bindings, name, value)
	return tc
}

func (tc TestContext) WithCluster(ctx context.Context, name string, cluster clusters.Cluster) TestContext {
	tc.clusters = tc.clusters.Register(name, cluster)
	return tc
}

func (tc TestContext) WithCurrentCluster(ctx context.Context, name string) TestContext {
	tc.cluster = tc.Cluster(name)
	return tc
}

func (tc TestContext) WithDryRun(ctx context.Context, dryRun bool) TestContext {
	tc.dryRun = dryRun
	return tc
}
