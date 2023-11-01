package logging

import (
	"fmt"

	"github.com/fatih/color"
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

type TestLogger struct {
	messages []string
}

func (tl *TestLogger) Log(args ...interface{}) {
	for _, arg := range args {
		tl.messages = append(tl.messages, fmt.Sprint(arg))
	}
}

func (tl *TestLogger) Helper() {}

type MockLogger struct {
	Logs []string
}

func (m *MockLogger) WithResource(resource ctrlclient.Object) Logger {
	return m
}

func (m *MockLogger) Log(operation string, color *color.Color, args ...interface{}) {
	message := fmt.Sprintf("%s: %v", operation, args)
	m.Logs = append(m.Logs, message)
}
