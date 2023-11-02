package logging

import (
	"context"

	"github.com/kyverno/kyverno/ext/output/color"
)

type contextKey struct{}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextKey{}).(Logger); ok {
		return v
	}
	return nil
}

func IntoContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func Log(ctx context.Context, operation string, color *color.Color, args ...interface{}) {
	logger := FromContext(ctx)
	if logger != nil {
		logger.Log(operation, color, args...)
	}
}
