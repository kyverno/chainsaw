package context

import (
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
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
	catch               []v1alpha1.CatchFinally
	cluster             clusters.Cluster
	clusters            clusters.Registry
	compilers           compilers.Compilers
	delayBeforeCleanup  *time.Duration
	deletionPropagation metav1.DeletionPropagation
	dryRun              bool
	failFast            bool
	fullName            bool
	skipDelete          bool
	templating          bool
	terminationGrace    *time.Duration
	timeouts            v1alpha1.DefaultTimeouts
}

func MakeContext(bindings apis.Bindings, registry clusters.Registry) TestContext {
	return TestContext{
		Summary: &model.Summary{},
		Report: &model.Report{
			Name:      "chainsaw-report",
			StartTime: time.Now(),
		},
		bindings:  bindings,
		clusters:  registry,
		compilers: apis.DefaultCompilers,
	}
}

func EmptyContext() TestContext {
	return MakeContext(apis.NewBindings(), clusters.NewRegistry(nil))
}

func (tc *TestContext) Bindings() apis.Bindings {
	return tc.bindings
}

func (tc *TestContext) Catch() []v1alpha1.CatchFinally {
	return tc.catch
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

func (tc *TestContext) DelayBeforeCleanup() *time.Duration {
	return tc.delayBeforeCleanup
}

func (tc *TestContext) DeletionPropagation() metav1.DeletionPropagation {
	return tc.deletionPropagation
}

func (tc *TestContext) DryRun() bool {
	return tc.dryRun
}

func (tc *TestContext) FailFast() bool {
	return tc.failFast
}

func (tc *TestContext) FullName() bool {
	return tc.fullName
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

func (tc *TestContext) Timeouts() v1alpha1.DefaultTimeouts {
	return tc.timeouts
}

func (tc TestContext) WithBinding(name string, value any) TestContext {
	tc.bindings = apibindings.RegisterBinding(tc.bindings, name, value)
	return tc
}

func (tc TestContext) WithCatch(catch ...v1alpha1.CatchFinally) TestContext {
	tc.catch = append(tc.catch, catch...)
	return tc
}

func (tc TestContext) WithDefaultCompiler(name string) TestContext {
	tc.compilers = tc.compilers.WithDefaultCompiler(name)
	return tc
}

func (tc TestContext) WithCluster(name string, cluster clusters.Cluster) TestContext {
	tc.clusters = tc.clusters.Register(name, cluster)
	return tc
}

func (tc TestContext) WithCurrentCluster(name string) TestContext {
	tc.cluster = tc.Cluster(name)
	return tc
}

func (tc TestContext) WithDelayBeforeCleanup(delayBeforeCleanup *time.Duration) TestContext {
	tc.delayBeforeCleanup = delayBeforeCleanup
	return tc
}

func (tc TestContext) WithDeletionPropagation(deletionPropagation metav1.DeletionPropagation) TestContext {
	tc.deletionPropagation = deletionPropagation
	return tc
}

func (tc TestContext) WithDryRun(dryRun bool) TestContext {
	tc.dryRun = dryRun
	return tc
}

func (tc TestContext) WithFailFast(failFast bool) TestContext {
	tc.failFast = failFast
	return tc
}

func (tc TestContext) WithFullName(fullName bool) TestContext {
	tc.fullName = fullName
	return tc
}

func (tc TestContext) WithSkipDelete(skipDelete bool) TestContext {
	tc.skipDelete = skipDelete
	return tc
}

func (tc TestContext) WithTemplating(templating bool) TestContext {
	tc.templating = templating
	return tc
}

func (tc TestContext) WithTerminationGrace(terminationGrace *time.Duration) TestContext {
	tc.terminationGrace = terminationGrace
	return tc
}

func (tc TestContext) WithTimeouts(timeouts v1alpha1.Timeouts) TestContext {
	if new := timeouts.Apply; new != nil {
		tc.timeouts.Apply = *new
	}
	if new := timeouts.Assert; new != nil {
		tc.timeouts.Assert = *new
	}
	if new := timeouts.Cleanup; new != nil {
		tc.timeouts.Cleanup = *new
	}
	if new := timeouts.Delete; new != nil {
		tc.timeouts.Delete = *new
	}
	if new := timeouts.Error; new != nil {
		tc.timeouts.Error = *new
	}
	if new := timeouts.Exec; new != nil {
		tc.timeouts.Exec = *new
	}
	return tc
}
