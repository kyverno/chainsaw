package context

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/client/dryrun"
	apibindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type TestContext struct {
	*model.Summary
	*model.Report
	bindings            apis.Bindings
	compilers           compilers.Compilers
	cluster             clusters.Cluster
	clusters            clusters.Registry
	delayBeforeCleanup  *time.Duration
	deletionPropagation metav1.DeletionPropagation
	dryRun              bool
	skipDelete          bool
	templating          bool
	terminationGrace    *time.Duration
}

func MakeContext(bindings apis.Bindings, registry clusters.Registry) TestContext {
	return TestContext{
		Summary: &model.Summary{},
		Report: &model.Report{
			Name:      "chainsaw-report",
			StartTime: time.Now(),
		},
		bindings:  bindings,
		compilers: apis.DefaultCompilers,
		clusters:  registry,
		cluster:   nil,
	}
}

func EmptyContext() TestContext {
	return MakeContext(apis.NewBindings(), clusters.NewRegistry(nil))
}

func (tc *TestContext) Bindings() apis.Bindings {
	return tc.bindings
}

func (tc *TestContext) Compilers() compilers.Compilers {
	return tc.compilers
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

func (tc *TestContext) DeletionPropagation() metav1.DeletionPropagation {
	return tc.deletionPropagation
}

func (tc *TestContext) DelayBeforeCleanup() *time.Duration {
	return tc.delayBeforeCleanup
}

func (tc *TestContext) SkipDelete() bool {
	return tc.skipDelete
}

func (tc *TestContext) Templating() bool {
	return tc.templating
}

func (tc *TestContext) TerminationGrace() *time.Duration {
	return tc.terminationGrace
}

func (tc TestContext) WithBinding(ctx context.Context, name string, value any) TestContext {
	tc.bindings = apibindings.RegisterBinding(ctx, tc.bindings, name, value)
	return tc
}

func (tc TestContext) WithDefaultCompiler(name string) TestContext {
	tc.compilers = tc.compilers.WithDefaultCompiler(name)
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

func (tc TestContext) WithDelayBeforeCleanup(ctx context.Context, delayBeforeCleanup *time.Duration) TestContext {
	tc.delayBeforeCleanup = delayBeforeCleanup
	return tc
}

func (tc TestContext) WithDeletionPropagation(ctx context.Context, deletionPropagation metav1.DeletionPropagation) TestContext {
	tc.deletionPropagation = deletionPropagation
	return tc
}

func (tc TestContext) WithSkipDelete(ctx context.Context, skipDelete bool) TestContext {
	tc.skipDelete = skipDelete
	return tc
}

func (tc TestContext) WithTemplating(ctx context.Context, templating bool) TestContext {
	tc.templating = templating
	return tc
}

func (tc TestContext) WithTerminationGrace(ctx context.Context, terminationGrace *time.Duration) TestContext {
	tc.terminationGrace = terminationGrace
	return tc
}
