package kubectl

import (
	"errors"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func getGroupVersionKind(mapper meta.RESTMapper, resource v1alpha1.ResourceReference) (schema.GroupVersionKind, error) {
	if resource.Resource != "" {
		gvr, gv := schema.ParseResourceArg(resource.Resource)
		if gvr == nil {
			gvr = &schema.GroupVersionResource{Group: gv.Group, Resource: gv.Resource}
		}
		return mapper.KindFor(*gvr)
	}
	if resource.Kind != "" {
		gv, err := schema.ParseGroupVersion(resource.APIVersion)
		if err != nil {
			return schema.GroupVersionKind{}, err
		}
		return gv.WithKind(resource.Kind), nil
	}
	return schema.GroupVersionKind{}, errors.New("failed to get GVK")
}

func getMapping(client client.Client, resource v1alpha1.ResourceReference) (*meta.RESTMapping, error) {
	mapper := client.RESTMapper()
	gvk, err := getGroupVersionKind(mapper, resource)
	if err != nil {
		return nil, err
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}

func mapResource(client client.Client, resource v1alpha1.ResourceReference) (string, meta.RESTScope, error) {
	mapping, err := getMapping(client, resource)
	if err != nil {
		return "", nil, err
	}
	parts := make([]string, 0, 3)
	parts = append(parts, mapping.Resource.Resource)
	if mapping.Resource.Group != "" {
		parts = append(parts, mapping.Resource.Version, mapping.Resource.Group)
	}
	return strings.Join(parts, "."), mapping.Scope, nil
}
