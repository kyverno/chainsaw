package logging

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/utils/clock"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operationLogger struct {
	t      *testing.T
	clock  clock.PassiveClock
	prefix string
}

func NewOperationLogger(t *testing.T, clock clock.PassiveClock, prefixes ...string) Logger {
	t.Helper()
	return &operationLogger{
		t:      t,
		clock:  clock,
		prefix: strings.Join(prefixes, " | "),
	}
}

func (l *operationLogger) Log(args ...interface{}) {
	Log(l.t, l.clock, l.prefix, args...)
}

func (l *operationLogger) Logf(format string, args ...interface{}) {
	Logf(l.t, l.clock, l.prefix, format, args...)
}

func (l *operationLogger) WithName(prefix string) Logger {
	return &operationLogger{
		t:      l.t,
		clock:  l.clock,
		prefix: strings.Join([]string{l.prefix, prefix}, " | "),
	}
}

func (l *operationLogger) WithResource(key ctrlclient.ObjectKey, obj ctrlclient.Object) Logger {
	gvk := obj.GetObjectKind().GroupVersionKind()
	name := client.Name(key)
	return &operationLogger{
		t:      l.t,
		clock:  l.clock,
		prefix: strings.Join([]string{l.prefix, fmt.Sprintf("%s/%s", gvk.GroupVersion(), gvk.Kind), name}, " | "),
	}
}
