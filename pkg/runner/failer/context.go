package failer

import (
	"context"
)

type contextKey struct{}

func FromContext(ctx context.Context) Failer {
	if ctx != nil {
		if v, ok := ctx.Value(contextKey{}).(Failer); ok {
			return v
		}
	}
	return nil
}

func IntoContext(ctx context.Context, failer Failer) context.Context {
	return context.WithValue(ctx, contextKey{}, failer)
}
