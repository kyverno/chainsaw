package context

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/model"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
)

func InitContext(config model.Configuration, defaultCluster *rest.Config, values any) (TestContext, error) {
	tc := EmptyContext(clock.RealClock{})
	// cleanup options
	tc = tc.WithSkipDelete(config.Cleanup.SkipDelete)
	if config.Cleanup.DelayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(&config.Cleanup.DelayBeforeCleanup.Duration)
	}
	// templating options
	tc = tc.WithTemplating(config.Templating.Enabled)
	if config.Templating.Compiler != nil {
		tc = tc.WithDefaultCompiler(string(*config.Templating.Compiler))
	}
	// discovery options
	tc = tc.WithFullName(config.Discovery.FullName)
	// execution options
	tc = tc.WithFailFast(config.Execution.FailFast)
	if config.Execution.ForceTerminationGracePeriod != nil {
		tc = tc.WithTerminationGrace(&config.Execution.ForceTerminationGracePeriod.Duration)
	}
	// deletion options
	tc = tc.WithDeletionPropagation(config.Deletion.Propagation)
	// error options
	tc = tc.WithCatch(config.Error.Catch...)
	// timeouts
	tc = tc.WithTimeouts(v1alpha1.Timeouts{
		Apply:   &config.Timeouts.Apply,
		Assert:  &config.Timeouts.Assert,
		Cleanup: &config.Timeouts.Cleanup,
		Delete:  &config.Timeouts.Delete,
		Error:   &config.Timeouts.Error,
		Exec:    &config.Timeouts.Exec,
	})
	// values
	tc = withValues(tc, values)
	// clusters
	tc = WithClusters(tc, "", config.Clusters)
	if defaultCluster != nil {
		cluster, err := clusters.NewClusterFromConfig(defaultCluster)
		if err != nil {
			return tc, err
		}
		tc = tc.WithCluster(clusters.DefaultClient, cluster)
		return WithCurrentCluster(tc, clusters.DefaultClient)
	}
	return tc, nil
}
