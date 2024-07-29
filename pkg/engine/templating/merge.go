package templating

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TemplateAndMerge(ctx context.Context, obj unstructured.Unstructured, bindings binding.Bindings, templates ...v1alpha1.Any) (unstructured.Unstructured, error) {
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
