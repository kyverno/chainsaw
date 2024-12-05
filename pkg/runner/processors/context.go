package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
}

func setupContext(ctx context.Context, tc engine.Context, data contextData) (engine.Context, error) {
	if len(data.catch) > 0 {
		tc = tc.WithCatch(ctx, data.catch...)
	}
	if data.dryRun != nil {
		tc = tc.WithDryRun(ctx, *data.dryRun)
	}
	if data.delayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(ctx, &data.delayBeforeCleanup.Duration)
	}
	if data.deletionPropagation != nil {
		tc = tc.WithDeletionPropagation(ctx, *data.deletionPropagation)
	}
	if data.skipDelete != nil {
		tc = tc.WithSkipDelete(ctx, *data.skipDelete)
	}
	if data.templating != nil {
		tc = tc.WithTemplating(ctx, *data.templating)
	}
	if data.terminationGrace != nil {
		tc = tc.WithTerminationGrace(ctx, &data.terminationGrace.Duration)
	}
	tc = engine.WithClusters(ctx, tc, data.basePath, data.clusters)
	if data.cluster != nil {
		if _tc, err := engine.WithCurrentCluster(ctx, tc, *data.cluster); err != nil {
			return tc, err
		} else {
			tc = _tc
		}
	}
	return tc, nil
}

func setupNamespace(ctx context.Context, tc engine.Context, data namespaceData) (engine.Context, *corev1.Namespace, error) {
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
		tc = engine.WithNamespace(ctx, tc, ns.GetName())
	}
	return tc, ns, nil
}

func setupBindings(ctx context.Context, tc engine.Context, bindings ...v1alpha1.Binding) (engine.Context, error) {
	if _tc, err := engine.WithBindings(ctx, tc, bindings...); err != nil {
		return tc, err
	} else {
		tc = _tc
	}
	return tc, nil
}

func setupContextAndBindings(ctx context.Context, tc engine.Context, data contextData, bindings ...v1alpha1.Binding) (engine.Context, error) {
	if tc, err := setupContext(ctx, tc, data); err != nil {
		return tc, err
	} else {
		return setupBindings(ctx, tc, bindings...)
	}
}
