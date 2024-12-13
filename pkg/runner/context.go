package runner

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func InitContext(config model.Configuration, defaultCluster *rest.Config, values any) (enginecontext.TestContext, error) {
	tc := enginecontext.EmptyContext()
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
	tc = enginecontext.WithValues(tc, values)
	// clusters
	tc = enginecontext.WithClusters(tc, "", config.Clusters)
	if defaultCluster != nil {
		cluster, err := clusters.NewClusterFromConfig(defaultCluster)
		if err != nil {
			return tc, err
		}
		tc = tc.WithCluster(clusters.DefaultClient, cluster)
		return enginecontext.WithCurrentCluster(tc, clusters.DefaultClient)
	}
	return tc, nil
}

type namespaceData struct {
	cleaner   cleaner.CleanerCollector
	compilers compilers.Compilers
	name      string
	template  *v1alpha1.Projection
}

type contextData struct {
	basePath            string
	catch               []v1alpha1.CatchFinally
	cluster             *string
	clusters            v1alpha1.Clusters
	delayBeforeCleanup  *metav1.Duration
	deletionPropagation *metav1.DeletionPropagation
	dryRun              *bool
	skipDelete          *bool
	templating          *bool
	terminationGrace    *metav1.Duration
	timeouts            *v1alpha1.Timeouts
}

func setupContext(tc enginecontext.TestContext, data contextData) (enginecontext.TestContext, error) {
	if len(data.catch) > 0 {
		tc = tc.WithCatch(data.catch...)
	}
	if data.dryRun != nil {
		tc = tc.WithDryRun(*data.dryRun)
	}
	if data.delayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(&data.delayBeforeCleanup.Duration)
	}
	if data.deletionPropagation != nil {
		tc = tc.WithDeletionPropagation(*data.deletionPropagation)
	}
	if data.skipDelete != nil {
		tc = tc.WithSkipDelete(*data.skipDelete)
	}
	if data.templating != nil {
		tc = tc.WithTemplating(*data.templating)
	}
	if data.terminationGrace != nil {
		tc = tc.WithTerminationGrace(&data.terminationGrace.Duration)
	}
	if data.timeouts != nil {
		tc = tc.WithTimeouts(*data.timeouts)
	}
	tc = enginecontext.WithClusters(tc, data.basePath, data.clusters)
	if data.cluster != nil {
		if _tc, err := enginecontext.WithCurrentCluster(tc, *data.cluster); err != nil {
			return tc, err
		} else {
			tc = _tc
		}
	}
	return tc, nil
}

func setupNamespace(ctx context.Context, tc enginecontext.TestContext, data namespaceData) (enginecontext.TestContext, *corev1.Namespace, error) {
	var ns *corev1.Namespace
	if namespace, err := buildNamespace(ctx, data.compilers, data.name, data.template, tc.Bindings()); err != nil {
		return tc, nil, err
	} else if _, clusterClient, err := tc.CurrentClusterClient(); err != nil {
		return tc, nil, err
	} else if clusterClient != nil {
		if err := clusterClient.Get(ctx, client.Key(namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				return tc, nil, err
			} else if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
				return tc, nil, err
			} else if data.cleaner != nil {
				data.cleaner.Add(clusterClient, namespace)
			}
		}
		ns = namespace
	}
	if ns != nil {
		tc = enginecontext.WithNamespace(tc, ns.GetName())
	}
	return tc, ns, nil
}

func setupBindings(tc enginecontext.TestContext, bindings ...v1alpha1.Binding) (enginecontext.TestContext, error) {
	if _tc, err := enginecontext.WithBindings(tc, bindings...); err != nil {
		return tc, err
	} else {
		tc = _tc
	}
	return tc, nil
}

func setupContextAndBindings(tc enginecontext.TestContext, data contextData, bindings ...v1alpha1.Binding) (enginecontext.TestContext, error) {
	if tc, err := setupContext(tc, data); err != nil {
		return tc, err
	} else {
		return setupBindings(tc, bindings...)
	}
}
