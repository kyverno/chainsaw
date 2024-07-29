package testing

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
)

// TODO: not thread safe
type FakeLogger struct {
	Logs     []string
	numCalls int
}

func (f *FakeLogger) WithResource(resource client.Object) Logger {
	defer func() { f.numCalls++ }()
	return f
}

func (f *FakeLogger) Log(operation Operation, status Status, color *color.Color, args ...fmt.Stringer) {
	defer func() { f.numCalls++ }()
	message := fmt.Sprintf("%s: %s - %v", operation, status, args)
	f.Logs = append(f.Logs, message)
}

func (f *FakeLogger) NumCalls() int {
	return f.numCalls
}
