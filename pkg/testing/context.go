package testing

import (
	"context"
	"testing"
	"time"
)

var MainStart = testing.MainStart

type (
	InternalTest = testing.InternalTest
	T            = testing.T
)

type contextKey struct{}

type tTest interface {
	Cleanup(func())
	Deadline() (deadline time.Time, ok bool)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Parallel()
	Run(name string, f func(t *T)) bool
	Setenv(key, value string)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
}

func FromContext(ctx context.Context) tTest {
	if v, ok := ctx.Value(contextKey{}).(tTest); ok {
		return v
	}
	return nil
}

func IntoContext(ctx context.Context, t tTest) context.Context {
	t.Helper()
	return context.WithValue(ctx, contextKey{}, t)
}
