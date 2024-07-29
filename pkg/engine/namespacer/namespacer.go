package namespacer

import (
	"errors"

	"github.com/kyverno/chainsaw/pkg/client"
)

type Namespacer interface {
	Apply(client.Client, client.Object) error
	GetNamespace() string
}

type namespacer struct {
	namespace string
}

func New(namespace string) Namespacer {
	return &namespacer{
		namespace: namespace,
	}
}

func (n *namespacer) Apply(client client.Client, resource client.Object) error {
	if resource == nil {
		return errors.New("resource is nil")
	}
	if resource.GetNamespace() == "" {
		namespaced, err := client.IsObjectNamespaced(resource)
		if err != nil {
			return err
		}
		if namespaced {
			resource.SetNamespace(n.namespace)
		}
	}
	return nil
}

func (n *namespacer) GetNamespace() string {
	return n.namespace
}
