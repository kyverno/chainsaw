package logging

import (
	"fmt"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/utils/clock"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const eraser = "\b\b\b\b\b\b\b\b\b"

type logger struct {
	t        *testing.T
	clock    clock.PassiveClock
	test     string
	step     string
	resource ctrlclient.Object
}

func NewLogger(t *testing.T, clock clock.PassiveClock, test string, step string) Logger {
	t.Helper()
	return &logger{
		t:     t,
		clock: clock,
		test:  test,
		step:  step,
	}
}

func (l *logger) Log(operation string, args ...interface{}) {
	a := make([]interface{}, 0, len(args)+1)
	if l.resource == nil {
		a = append(a, fmt.Sprintf("%s%s | %s | %s | %s |", eraser, l.clock.Now().Format("15:04:05"), operation, l.test, l.step))
	} else {
		gvk := l.resource.GetObjectKind().GroupVersionKind()
		name := client.Name(client.ObjectKey(l.resource))
		a = append(a, fmt.Sprintf("%s%s | %s | %s | %s | %s | %s |", eraser, l.clock.Now().Format("15:04:05"), operation, l.test, l.step, fmt.Sprintf("%s/%s", gvk.GroupVersion(), gvk.Kind), name))
	}
	a = append(a, args...)
	l.t.Log(a...)
}

func (l *logger) WithResource(resource ctrlclient.Object) Logger {
	return &logger{
		t:        l.t,
		clock:    l.clock,
		test:     l.test,
		step:     l.step,
		resource: resource,
	}
}
