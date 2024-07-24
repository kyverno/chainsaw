package logging

import (
	"context"
)

type contextKey struct{}

func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		if v, ok := ctx.Value(contextKey{}).(Logger); ok {
			return v
		}
	}
	return nil
}

func IntoContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}
