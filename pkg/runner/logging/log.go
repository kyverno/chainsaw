package logging

import (
	"fmt"
	"testing"
	"time"
)

func Log(t *testing.T, prefix string, args ...interface{}) {
	t.Helper()
	a := make([]interface{}, 0, len(args)+1)
	a = append(a, fmt.Sprintf("%s | %s |", time.Now().Format("15:04:05"), prefix))
	a = append(a, args...)
	t.Log(a...)
}

func Logf(t *testing.T, prefix string, format string, args ...interface{}) {
	t.Helper()
	Log(t, prefix, fmt.Sprintf(format, args...))
}
