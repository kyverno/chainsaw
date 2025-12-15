package logging

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
)

type Sink interface {
	Log(string, string, string, Operation, Status, client.Object, *color.Color, ...fmt.Stringer)
}

type SinkFunc func(string, string, string, Operation, Status, client.Object, *color.Color, ...fmt.Stringer)

func (f SinkFunc) Log(test string, scenario string, step string, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
	f(test, scenario, step, operation, status, obj, color, args...)
}
