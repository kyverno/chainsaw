package testing

import (
	"fmt"

	"github.com/fatih/color"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: not thread safe
type FakeLogger struct {
	Logs     []string
	numCalls int
}

func (f *FakeLogger) WithResource(resource ctrlclient.Object) Logger {
	defer func() { f.numCalls++ }()
	return f
}

func (f *FakeLogger) Log(operation string, color *color.Color, args ...interface{}) {
	defer func() { f.numCalls++ }()
	message := fmt.Sprintf("%s: %v", operation, args)
	f.Logs = append(f.Logs, message)
}

func (f *FakeLogger) NumCalls() int {
	return f.numCalls
}
