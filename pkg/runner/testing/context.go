package testing

import (
	"context"
	"testing"
)

var MainStart = testing.MainStart

type (
	InternalTest = testing.InternalTest
	T            = testing.T
)

type contextKey struct{}

func FromContext(ctx context.Context) *testing.T {
	if v, ok := ctx.Value(contextKey{}).(*testing.T); ok {
		return v
	}
	return nil
}

func IntoContext(ctx context.Context, t *testing.T) context.Context {
	t.Helper()
	return context.WithValue(ctx, contextKey{}, t)
}
