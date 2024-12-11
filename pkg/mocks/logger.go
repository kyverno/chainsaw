package mocks

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/logging"
)

type Logger struct {
	Logs     []string
	numCalls int
}

func (f *Logger) Log(_ context.Context, operation logging.Operation, status logging.Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
	defer func() { f.numCalls++ }()
	message := fmt.Sprintf("%s: %s - %v", operation, status, args)
	f.Logs = append(f.Logs, message)
}

func (f *Logger) NumCalls() int {
	return f.numCalls
}
