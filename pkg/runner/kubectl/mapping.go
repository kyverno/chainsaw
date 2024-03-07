package kubectl

import (
	"errors"
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func mapResource(client client.Client, bindings binding.Bindings, resource v1alpha1.ResourceReference) (string, bool, error) {
	if resource.Resource != "" {
		if resource, err := ConvertString(resource.Resource, bindings); err != nil {
			return "", false, err
		} else {
			return mapResourceFromResource(client, resource)
		}
	}
	if resource.APIVersion != "" && resource.Kind != "" {
		if apiVersion, err := ConvertString(resource.APIVersion, bindings); err != nil {
			return "", false, err
		} else if kind, err := ConvertString(resource.Kind, bindings); err != nil {
			return "", false, err
		} else {
			return mapResourceFromKind(client, apiVersion, kind)
		}
	}
	return "", false, errors.New("failed to map resource, either kind or resource must be specified")
}

func mapResourceFromResource(client client.Client, resource string) (string, bool, error) {
	gvr, gv := schema.ParseResourceArg(resource)
	if gvr == nil {
		gvr = &schema.GroupVersionResource{Group: gv.Group, Resource: gv.Resource}
	}
	mapper := client.RESTMapper()
	gvk, err := mapper.KindFor(*gvr)
	// if we have an error, it may be because the resource name is a short one
	if err != nil {
		return resource, false, nil
	}
	return mapResourceFromGVK(mapper, gvk)
}

func mapResourceFromKind(client client.Client, apiVersion string, kind string) (string, bool, error) {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return "", false, err
	}
	return mapResourceFromGVK(client.RESTMapper(), gv.WithKind(kind))
}

func mapResourceFromGVK(mapper meta.RESTMapper, gvk schema.GroupVersionKind) (string, bool, error) {
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", false, err
	}
	clustered := mapping.Scope.Name() == meta.RESTScopeNameRoot
	if mapping.Resource.Group == "" {
		return mapping.Resource.Resource, clustered, nil
	}
	parts := []string{
		mapping.Resource.Resource,
		mapping.Resource.Version,
		mapping.Resource.Group,
	}
	return strings.Join(parts, "."), clustered, nil
}
