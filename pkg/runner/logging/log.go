package logging

import (
	"fmt"
	"testing"
	"time"
)

func Log(t *testing.T, prefix string, args ...interface{}) {
	t.Helper()
	args = append([]interface{}{
		fmt.Sprintf("%s | %s |", time.Now().Format("15:04:05"), prefix),
	}, args...)
	t.Log(args...)
}

func Logf(t *testing.T, prefix string, format string, args ...interface{}) {
	t.Helper()
	Log(t, prefix, fmt.Sprintf(format, args...))
}
