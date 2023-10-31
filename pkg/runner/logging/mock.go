package logging

import (
	"time"

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

// This is a mock clock for testing purposes
type mockClock struct {
	time time.Time
}

func (m *mockClock) Now() time.Time {
	return m.time
}

func (m *mockClock) Since(time.Time) time.Duration {
	return 0
}
