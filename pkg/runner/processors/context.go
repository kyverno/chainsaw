package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type namespaceData struct {
	name     string
	template *v1alpha1.Any
	cleaner  cleaner.CleanerCollector
}

type contextData struct {
	basePath  string
	bindings  []v1alpha1.Binding
	cluster   *string
	clusters  v1alpha1.Clusters
	dryRun    *bool
	namespace *namespaceData
}

func setupContextData(ctx context.Context, tc engine.Context, data contextData) (engine.Context, *corev1.Namespace, error) {
	tc = engine.WithClusters(ctx, tc, data.basePath, data.clusters)
	if data.dryRun != nil {
		tc = tc.WithDryRun(ctx, *data.dryRun)
	}
	if data.cluster != nil {
		if _tc, err := engine.WithCurrentCluster(ctx, tc, *data.cluster); err != nil {
			return tc, nil, err
		} else {
			tc = _tc
		}
	}
	var ns *corev1.Namespace
	if data.namespace != nil {
		if namespace, err := buildNamespace(ctx, data.namespace.name, data.namespace.template, tc.Bindings()); err != nil {
			return tc, nil, err
		} else if _, clusterClient, err := tc.CurrentClusterClient(); err != nil {
			return tc, nil, err
		} else {
			if err := clusterClient.Get(ctx, client.Key(namespace), namespace.DeepCopy()); err != nil {
				if !errors.IsNotFound(err) {
					return tc, nil, err
				} else if err := clusterClient.Create(ctx, namespace.DeepCopy()); err != nil {
					return tc, nil, err
				} else if data.namespace.cleaner != nil {
					data.namespace.cleaner.Add(clusterClient, namespace)
				}
			}
			ns = namespace
		}
		tc = engine.WithNamespace(ctx, tc, ns.GetName())
	}
	if _tc, err := engine.WithBindings(ctx, tc, data.bindings...); err != nil {
		return tc, ns, err
	} else {
		tc = _tc
	}
	return tc, ns, nil
}
