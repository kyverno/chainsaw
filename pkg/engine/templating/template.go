package templating

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/mutate"
	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Template(ctx context.Context, compilers compilers.Compilers, tpl v1alpha1.Projection, value any, bindings apis.Bindings) (any, error) {
	return mutate.Mutate(ctx, nil, mutate.Parse(ctx, tpl.Value()), value, bindings, compilers)
}

func TemplateAndMerge(ctx context.Context, compilers compilers.Compilers, obj unstructured.Unstructured, bindings apis.Bindings, templates ...v1alpha1.Projection) (unstructured.Unstructured, error) {
	for _, modifier := range templates {
		patch, err := Template(ctx, compilers, modifier, obj.UnstructuredContent(), bindings)
		if err != nil {
			return obj, err
		}
		patched := mapsutils.Merge(obj.UnstructuredContent(), convertMap(patch))
		obj.SetUnstructuredContent(patched)
	}
	return obj, nil
}
