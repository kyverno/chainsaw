package mutate

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Merge(ctx context.Context, obj unstructured.Unstructured, bindings binding.Bindings, modifiers ...v1alpha1.Any) (unstructured.Unstructured, error) {
	for _, modifier := range modifiers {
		patch, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, modifier.Value), obj.UnstructuredContent(), bindings, template.WithFunctionCaller(functions.Caller))
		if err != nil {
			return obj, err
		}
		obj.SetUnstructuredContent(mapsutils.Merge(obj.UnstructuredContent(), convertMap(patch)))
	}
	return obj, nil
}
