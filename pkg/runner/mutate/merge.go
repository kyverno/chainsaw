package mutate

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	mapsutils "github.com/kyverno/chainsaw/pkg/utils/maps"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Merge(ctx context.Context, obj unstructured.Unstructured, bindings binding.Bindings, modifiers ...v1alpha1.Modifier) (unstructured.Unstructured, error) {
	for _, modifier := range modifiers {
		// if a match is specified, skip the modifier if the resource doesn't match
		if modifier.Match != nil && modifier.Match.Value != nil {
			if errs, err := check.Check(ctx, obj.UnstructuredContent(), nil, modifier.Match); err != nil {
				return obj, err
			} else if len(errs) != 0 {
				continue
			}
		}
		merge := modifier.Merge
		if modifier.Label != nil {
			merge = &v1alpha1.Any{
				Value: map[string]any{
					"metadata": map[string]any{
						"labels": modifier.Label.Value,
					},
				},
			}
		} else if modifier.Annotate != nil {
			merge = &v1alpha1.Any{
				Value: map[string]any{
					"metadata": map[string]any{
						"annotations": modifier.Annotate.Value,
					},
				},
			}
		}
		if merge != nil {
			patch, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, merge.Value), obj.UnstructuredContent(), bindings, template.WithFunctionCaller(functions.Caller))
			if err != nil {
				return obj, err
			}
			obj.SetUnstructuredContent(mapsutils.Merge(obj.UnstructuredContent(), convertMap(patch)))
		}
	}
	return obj, nil
}
