package templating

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/mutate"
	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Template(ctx context.Context, tpl v1alpha1.Projection, value any, bindings binding.Bindings) (any, error) {
	return mutate.Mutate(ctx, nil, mutate.Parse(ctx, tpl.Value()), value, bindings, apis.DefaultCompilers.Jp.Options()...)
}

func TemplateAndMerge(ctx context.Context, obj unstructured.Unstructured, bindings binding.Bindings, templates ...v1alpha1.Projection) (unstructured.Unstructured, error) {
	for _, modifier := range templates {
		patch, err := Template(ctx, modifier, obj.UnstructuredContent(), bindings)
		if err != nil {
			return obj, err
		}
		patched := mapsutils.Merge(obj.UnstructuredContent(), convertMap(patch))
		obj.SetUnstructuredContent(patched)
	}
	return obj, nil
}
