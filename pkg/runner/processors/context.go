package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/pkg/ext/output/color"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceData struct {
	Cleaner   cleaner.CleanerCollector
	Compilers compilers.Compilers
	Name      string
	Template  *v1alpha1.Projection
}

type ContextData struct {
	BasePath            string
	Catch               []v1alpha1.CatchFinally
	Cluster             *string
	Clusters            v1alpha1.Clusters
	DelayBeforeCleanup  *metav1.Duration
	DeletionPropagation *metav1.DeletionPropagation
	DryRun              *bool
	SkipDelete          *bool
	Templating          *bool
	TerminationGrace    *metav1.Duration
	Timeouts            *v1alpha1.Timeouts
}

func SetupContext(ctx context.Context, tc engine.Context, data ContextData) (engine.Context, error) {
	if len(data.Catch) > 0 {
		tc = tc.WithCatch(ctx, data.Catch...)
	}
	if data.DryRun != nil {
		tc = tc.WithDryRun(ctx, *data.DryRun)
	}
	if data.DelayBeforeCleanup != nil {
		tc = tc.WithDelayBeforeCleanup(ctx, &data.DelayBeforeCleanup.Duration)
	}
	if data.DeletionPropagation != nil {
		tc = tc.WithDeletionPropagation(ctx, *data.DeletionPropagation)
	}
	if data.SkipDelete != nil {
		tc = tc.WithSkipDelete(ctx, *data.SkipDelete)
	}
	if data.Templating != nil {
		tc = tc.WithTemplating(ctx, *data.Templating)
	}
	if data.TerminationGrace != nil {
		tc = tc.WithTerminationGrace(ctx, &data.TerminationGrace.Duration)
	}
	if data.Timeouts != nil {
		tc = tc.WithTimeouts(ctx, *data.Timeouts)
	}
	tc = engine.WithClusters(ctx, tc, data.BasePath, data.Clusters)
	if data.Cluster != nil {
		if _tc, err := engine.WithCurrentCluster(ctx, tc, *data.Cluster); err != nil {
			return tc, err
		} else {
			tc = _tc
		}
	}
	return tc, nil
}

func SetupNamespace(ctx context.Context, tc engine.Context, data NamespaceData) (engine.Context, *corev1.Namespace, error) {
	var ns *corev1.Namespace
	if namespace, err := buildNamespace(ctx, data.Compilers, data.Name, data.Template, tc.Bindings()); err != nil {
		return tc, nil, err
	} else if _, clusterClient, err := tc.CurrentClusterClient(); err != nil {
		return tc, nil, err
	} else if clusterClient != nil {
		if err := clusterClient.Get(ctx, client.Key(namespace), namespace.DeepCopy()); err != nil {
			if !errors.IsNotFound(err) {
				return tc, nil, err
			} else if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
				return tc, nil, err
			} else if data.Cleaner != nil {
				data.Cleaner.Add(clusterClient, namespace)
			}
		}
		ns = namespace
	}
	if ns != nil {
		tc = engine.WithNamespace(ctx, tc, ns.GetName())
	}
	return tc, ns, nil
}

func SetupBindings(ctx context.Context, tc engine.Context, bindings ...v1alpha1.Binding) (engine.Context, error) {
	if _tc, err := engine.WithBindings(ctx, tc, bindings...); err != nil {
		return tc, err
	} else {
		tc = _tc
	}
	return tc, nil
}

func SetupCleanup(ctx context.Context, t testing.TTest, failer failer.Failer, tc engine.Context) cleaner.CleanerCollector {
	if tc.SkipDelete() {
		return nil
	}
	cleaner := cleaner.New(tc.Timeouts().Cleanup.Duration, nil, tc.DeletionPropagation())
	t.Cleanup(func() {
		if !cleaner.Empty() {
			logging.Log(ctx, logging.Cleanup, logging.BeginStatus, color.BoldFgCyan)
			defer func() {
				logging.Log(ctx, logging.Cleanup, logging.EndStatus, color.BoldFgCyan)
			}()
			for _, err := range cleaner.Run(ctx, nil) {
				logging.Log(ctx, logging.Cleanup, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
				failer.Fail(ctx, t)
			}
		}
	})
	return cleaner
}

func setupContextAndBindings(ctx context.Context, tc engine.Context, data ContextData, bindings ...v1alpha1.Binding) (engine.Context, error) {
	if tc, err := SetupContext(ctx, tc, data); err != nil {
		return tc, err
	} else {
		return SetupBindings(ctx, tc, bindings...)
	}
}
