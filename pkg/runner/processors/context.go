package processors

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/model"
)

type contextData struct {
	basePath string
	clusters v1alpha1.Clusters
	cluster  *string
	bindings []v1alpha1.Binding
}

func setupContextData(ctx context.Context, tc model.TestContext, data contextData) (model.TestContext, error) {
	tc = model.WithClusters(ctx, tc, data.basePath, data.clusters)
	if data.cluster != nil {
		if _, _, _tc, err := model.WithCurrentCluster(ctx, tc, *data.cluster); err != nil {
			return tc, err
		} else {
			tc = _tc
		}
	}
	if _tc, err := model.WithBindings(ctx, tc, data.bindings...); err != nil {
		return tc, err
	} else {
		tc = _tc
	}
	return tc, nil
}
