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
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Parallel()
	Run(name string, f func(t *T)) bool
	Setenv(key, value string)
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
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
