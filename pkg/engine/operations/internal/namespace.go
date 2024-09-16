package internal

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
)

func ApplyNamespacer(namespacer namespacer.Namespacer, client client.Client, obj client.Object) error {
	if namespacer == nil {
		return nil
	}
	return namespacer.Apply(client, obj)
}
