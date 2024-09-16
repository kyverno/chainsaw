package testing

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type FakeLoader struct {
	LoadFn   func(int, []byte) (schema.GroupVersionKind, unstructured.Unstructured, error)
	numCalls int
}

func (f *FakeLoader) Load(data []byte) (schema.GroupVersionKind, unstructured.Unstructured, error) {
	defer func() { f.numCalls++ }()
	return f.LoadFn(f.numCalls, data)
}
