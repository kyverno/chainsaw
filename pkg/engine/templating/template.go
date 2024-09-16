package templating

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/chainsaw/pkg/mutate"
	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Template(ctx context.Context, tpl v1alpha1.Any, value any, bindings binding.Bindings) (any, error) {
	return mutate.Mutate(ctx, nil, mutate.Parse(ctx, tpl.Value), value, bindings, template.WithFunctionCaller(functions.Caller()))
}

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
