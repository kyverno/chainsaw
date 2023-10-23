package logging

import (
	"testing"

	"k8s.io/utils/clock"
)

type testLogger struct {
	t      *testing.T
	clock  clock.PassiveClock
	prefix string
}

func NewTestLogger(t *testing.T, clock clock.PassiveClock) Logger {
	t.Helper()
	prefix := t.Name()
	return &testLogger{
		t:      t,
		clock:  clock,
		prefix: prefix,
	}
}

func (l *testLogger) Log(args ...interface{}) {
	l.t.Helper()
	Log(l.t, l.clock, l.prefix, args...)
}

func (l *testLogger) Logf(format string, args ...interface{}) {
	l.t.Helper()
	Logf(l.t, l.clock, l.prefix, format, args...)
}
