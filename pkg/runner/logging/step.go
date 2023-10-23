package logging

import (
	"path/filepath"
	"testing"

	"k8s.io/utils/clock"
)

type stepLogger struct {
	t      *testing.T
	clock  clock.PassiveClock
	prefix string
}

func NewStepLogger(t *testing.T, clock clock.PassiveClock, step string) Logger {
	t.Helper()
	prefix := filepath.Join(t.Name(), step)
	return &stepLogger{
		t:      t,
		clock:  clock,
		prefix: prefix,
	}
}

func (l *stepLogger) Log(args ...interface{}) {
	l.t.Helper()
	Log(l.t, l.clock, l.prefix, args...)
}

func (l *stepLogger) Logf(format string, args ...interface{}) {
	l.t.Helper()
	Logf(l.t, l.clock, l.prefix, format, args...)
}
