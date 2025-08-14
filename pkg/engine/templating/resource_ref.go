package templating

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/utils/maps"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ResourceRef(ctx context.Context, c compilers.Compilers, obj *unstructured.Unstructured, bindings apis.Bindings) error {
	if obj == nil {
		return nil
	}
	apiVersion := obj.GetAPIVersion()
	kind := obj.GetKind()
	// this is not a valid resource (non resource assertion maybe ?)
	if apiVersion == "" || kind == "" {
		return nil
	}
	var temp unstructured.Unstructured
	temp.SetAPIVersion(apiVersion)
	temp.SetKind(kind)
	temp.SetName(obj.GetName())
	temp.SetNamespace(obj.GetNamespace())
	temp.SetLabels(getStringLabels(obj))
	template := v1alpha1.NewProjection(temp.UnstructuredContent())
	if merged, err := TemplateAndMerge(ctx, c, temp, bindings, template); err != nil {
		return err
	} else {
		obj.Object = maps.Merge(obj.Object, merged.UnstructuredContent())
	}
	return nil
}

func getStringLabels(obj *unstructured.Unstructured) map[string]string {
	if labels, ok, _ := unstructured.NestedFieldNoCopy(obj.UnstructuredContent(), "metadata", "labels"); ok {
		if labelsMap, ok := labels.(map[string]interface{}); ok {
			tempLabels := make(map[string]string)
			for k, v := range labelsMap {
				if stringValue, ok := v.(string); ok {
					tempLabels[k] = stringValue
				}
			}
			return tempLabels
		}
	}
	return nil
}
