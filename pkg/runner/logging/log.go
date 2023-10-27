package logging

import (
	"fmt"
	"testing"

	"k8s.io/utils/clock"
)

func Log(t *testing.T, clock clock.PassiveClock, prefix string, args ...interface{}) {
	t.Helper()
	a := make([]interface{}, 0, len(args)+1)
	a = append(a, fmt.Sprintf("\b\b\b\b\b\b\b\b\b%s | %s |", clock.Now().Format("15:04:05"), prefix))
	a = append(a, args...)
	t.Log(a...)
}

func Logf(t *testing.T, clock clock.PassiveClock, prefix string, format string, args ...interface{}) {
	t.Helper()
	Log(t, clock, prefix, fmt.Sprintf(format, args...))
}
