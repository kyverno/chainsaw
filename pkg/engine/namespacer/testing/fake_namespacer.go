package testing

import (
	"github.com/kyverno/chainsaw/pkg/client"
)

type FakeNamespacer struct {
	ApplyFn        func(call int, client client.Client, obj client.Object) error
	GetNamespaceFn func(call int) string
	numCalls       int
}

func (n *FakeNamespacer) Apply(client client.Client, obj client.Object) error {
	defer func() { n.numCalls++ }()
	return n.ApplyFn(n.numCalls, client, obj)
}

func (n *FakeNamespacer) GetNamespace() string {
	defer func() { n.numCalls++ }()
	return n.GetNamespaceFn(n.numCalls)
}
