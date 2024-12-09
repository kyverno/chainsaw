package testing

import (
	"testing"
	"time"
)

var MainStart = testing.MainStart

type (
	InternalTest = testing.InternalTest
	T            = testing.T
)

type TTest interface {
	Cleanup(func())
	Deadline() (deadline time.Time, ok bool)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
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
	Skipf(format string, args ...any)
	Skipped() bool
	TempDir() string
}
