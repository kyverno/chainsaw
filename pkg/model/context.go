package model

import (
	"context"
	"path/filepath"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"k8s.io/client-go/rest"
)

type Timeouts struct {
	Apply   time.Duration
	Assert  time.Duration
	Cleanup time.Duration
	Delete  time.Duration
	Error   time.Duration
	Exec    time.Duration
}

type TestContext struct {
	Summary
	timeouts Timeouts
	cleanup  bool
	bindings binding.Bindings
	clusters clusters.Registry
	cluster  string
}

func MakeContext(bindings binding.Bindings, registry clusters.Registry) TestContext {
	return TestContext{
		Summary:  &summary{},
		cleanup:  true,
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

func (tc *TestContext) Cleanup() bool {
	return tc.cleanup
}

func (tc *TestContext) Clusters() clusters.Registry {
	return tc.clusters
}

func (tc *TestContext) Cluster() (*rest.Config, client.Client, error) {
	return tc.clusters.Resolve(false, tc.cluster)
}

func (tc *TestContext) Timeouts() Timeouts {
	return tc.timeouts
}

func (tc TestContext) WithBinding(ctx context.Context, name string, value any) TestContext {
	tc.bindings = apibindings.RegisterNamedBinding(ctx, tc.bindings, name, value)
	return tc
}

func (tc TestContext) WithCleanup(ctx context.Context, cleanup bool) TestContext {
	tc.cleanup = cleanup
	return tc
}

func (tc TestContext) WithCluster(ctx context.Context, name string, cluster clusters.Cluster) TestContext {
	tc.clusters = tc.clusters.Register(name, cluster)
	return tc
}

func (tc TestContext) WithTimeouts(ctx context.Context, timeouts Timeouts) TestContext {
	tc.timeouts = timeouts
	return tc
}

func WithBindings(ctx context.Context, tc TestContext, variables ...v1alpha1.Binding) (TestContext, error) {
	bindings, err := apibindings.RegisterBindings(ctx, tc.Bindings(), variables...)
	if err != nil {
		return tc, err
	}
	tc.bindings = bindings
	return tc, nil
}

func WithClusters(ctx context.Context, tc TestContext, basePath string, c map[string]v1alpha1.Cluster) TestContext {
	for name, cluster := range c {
		kubeconfig := filepath.Join(basePath, cluster.Kubeconfig)
		cluster := clusters.NewClusterFromKubeconfig(kubeconfig, cluster.Context)
		tc = tc.WithCluster(ctx, name, cluster)
	}
	return tc
}

func WithDefaultTimeouts(ctx context.Context, tc TestContext, timeouts v1alpha1.DefaultTimeouts) TestContext {
	return tc.WithTimeouts(ctx, Timeouts{
		Apply:   timeouts.Apply.Duration,
		Assert:  timeouts.Assert.Duration,
		Cleanup: timeouts.Cleanup.Duration,
		Delete:  timeouts.Delete.Duration,
		Error:   timeouts.Error.Duration,
		Exec:    timeouts.Exec.Duration,
	})
}

func WithTimeouts(ctx context.Context, tc TestContext, timeouts v1alpha1.Timeouts) TestContext {
	old := tc.timeouts
	if new := timeouts.Apply; new != nil {
		old.Apply = new.Duration
	}
	if new := timeouts.Assert; new != nil {
		old.Assert = new.Duration
	}
	if new := timeouts.Cleanup; new != nil {
		old.Cleanup = new.Duration
	}
	if new := timeouts.Delete; new != nil {
		old.Delete = new.Duration
	}
	if new := timeouts.Error; new != nil {
		old.Error = new.Duration
	}
	if new := timeouts.Exec; new != nil {
		old.Exec = new.Duration
	}
	return tc.WithTimeouts(ctx, old)
}

func WithValues(ctx context.Context, tc TestContext, values any) TestContext {
	return tc.WithBinding(ctx, "values", values)
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
