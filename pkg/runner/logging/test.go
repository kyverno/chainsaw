package logging

import (
	"testing"
)

type testLogger struct {
	t      *testing.T
	prefix string
}

func NewTestLogger(t *testing.T) Logger {
	t.Helper()
	prefix := t.Name()
	return &testLogger{
		t:      t,
		prefix: prefix,
	}
}

func (l *testLogger) Log(args ...interface{}) {
	l.t.Helper()
	Log(l.t, l.prefix, args...)
}

func (l *testLogger) Logf(format string, args ...interface{}) {
	l.t.Helper()
	Logf(l.t, l.prefix, format, args...)
}
