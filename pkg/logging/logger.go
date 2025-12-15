package logging

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
)

type Logger interface {
	Log(context.Context, Operation, Status, client.Object, *color.Color, ...fmt.Stringer)
}

type LoggerFunc func(context.Context, Operation, Status, client.Object, *color.Color, ...fmt.Stringer)

func (f LoggerFunc) Log(ctx context.Context, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
	f(ctx, operation, status, obj, color, args...)
}

func NewLogger(test, scenario, step string) LoggerFunc {
	return func(ctx context.Context, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
		if sink := getSink(ctx); sink != nil {
			sink.Log(test, scenario, step, operation, status, obj, color, args...)
		}
	}
}
