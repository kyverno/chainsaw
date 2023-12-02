package internal

import (
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ApplyNamespacer(namespacer namespacer.Namespacer, obj client.Object) error {
	if namespacer == nil {
		return nil
	}
	return namespacer.Apply(obj)
}
