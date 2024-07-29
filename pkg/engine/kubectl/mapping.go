package kubectl

import (
	"context"
	"errors"
	"strings"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func mapResource(ctx context.Context, client client.Client, tc binding.Bindings, resource v1alpha1.ObjectType) (string, bool, error) {
	if resource.APIVersion != "" && resource.Kind != "" {
		if apiVersion, err := templating.String(ctx, resource.APIVersion, tc); err != nil {
			return "", false, err
		} else if kind, err := templating.String(ctx, resource.Kind, tc); err != nil {
			return "", false, err
		} else {
			return mapResourceFromApiVersionAndKind(client, apiVersion, kind)
		}
	}
	return "", false, errors.New("failed to map resource, either kind or resource must be specified")
}

func mapResourceFromApiVersionAndKind(client client.Client, apiVersion string, kind string) (string, bool, error) {
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
