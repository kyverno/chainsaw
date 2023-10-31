package logging

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// test resource for testing
type testResource struct {
	ctrlclient.Object
	name      string
	namespace string
	gvk       schema.GroupVersionKind
}

func (f *testResource) GetName() string {
	return f.name
}

func (f *testResource) GetNamespace() string {
	return f.namespace
}

func (f *testResource) GetObjectKind() schema.ObjectKind {
	return &objectKind{
		gvk: f.gvk,
	}
}

type objectKind struct {
	gvk schema.GroupVersionKind
}

func (o *objectKind) GroupVersionKind() schema.GroupVersionKind {
	return o.gvk
}

func (o *objectKind) SetGroupVersionKind(kind schema.GroupVersionKind) {
	o.gvk = kind
}
