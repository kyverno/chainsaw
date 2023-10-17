package logging

import (
	"path/filepath"
	"testing"
)

type stepLogger struct {
	t      *testing.T
	prefix string
}

func NewStepLogger(t *testing.T, step string) Logger {
	t.Helper()
	prefix := filepath.Join(t.Name(), step)
	return &stepLogger{
		t:      t,
		prefix: prefix,
	}
}

func (l *stepLogger) Log(args ...interface{}) {
	l.t.Helper()
	Log(l.t, l.prefix, args...)
}

func (l *stepLogger) Logf(format string, args ...interface{}) {
	l.t.Helper()
	Logf(l.t, l.prefix, format, args...)
}
