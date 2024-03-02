package kubectl

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func getMapping(client client.Client, resource string) (*meta.RESTMapping, error) {
	mapper := client.RESTMapper()
	gvr, gv := schema.ParseResourceArg(resource)
	if gvr == nil {
		gvr = &schema.GroupVersionResource{Group: gv.Group, Resource: gv.Resource}
	}
	gvk, err := mapper.KindFor(*gvr)
	if err != nil {
		return nil, err
	}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}
	return mapping, nil
}
