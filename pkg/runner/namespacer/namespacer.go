package namespacer

import (
	"github.com/kyverno/chainsaw/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Namespacer interface {
	Apply(ctrlclient.Object) error
}

type namespacer struct {
	c         client.Client
	namespace string
}

func New(c client.Client, namespace string) Namespacer {
	return &namespacer{
		c:         c,
		namespace: namespace,
	}
}

func (n *namespacer) Apply(resource ctrlclient.Object) error {
	if resource.GetNamespace() == "" {
		namespaced, err := n.c.IsObjectNamespaced(resource)
		if err != nil {
			return err
		}
		if namespaced {
			resource.SetNamespace(n.namespace)
		}
	}
	return nil
}
