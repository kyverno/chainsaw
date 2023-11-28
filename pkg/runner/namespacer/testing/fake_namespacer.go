package testing

import (
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type FakeNamespacer struct {
	ApplyFn        func(obj ctrlclient.Object, call int) error
	GetNamespaceFn func(call int) string
	numCalls       int
}

func (n *FakeNamespacer) Apply(obj ctrlclient.Object) error {
	defer func() { n.numCalls++ }()
	return n.ApplyFn(obj, n.numCalls)
}

func (n *FakeNamespacer) GetNamespace() string {
	defer func() { n.numCalls++ }()
	return n.GetNamespaceFn(n.numCalls)
}
