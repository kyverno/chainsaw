package logging

import (
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/utils/clock"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const eraser = "\b\b\b\b\b\b\b\b\b"

type tLogger interface {
	Log(args ...interface{})
	Helper()
}

type logger struct {
	t        tLogger
	clock    clock.PassiveClock
	test     string
	step     string
	resource ctrlclient.Object
}

func NewLogger(t tLogger, clock clock.PassiveClock, test string, step string) Logger {
	t.Helper()
	return &logger{
		t:     t,
		clock: clock,
		test:  test,
		step:  step,
	}
}

func (l *logger) Log(operation string, color *color.Color, args ...interface{}) {
	sprint := fmt.Sprint
	if color != nil {
		sprint = color.Sprint
	}
	a := make([]interface{}, 0, len(args)+1)
	prefix := fmt.Sprintf("%s%s | %s | %s | %s |", eraser, l.clock.Now().Format("15:04:05"), sprint(l.test), sprint(l.step), sprint(operation))
	if l.resource != nil {
		gvk := l.resource.GetObjectKind().GroupVersionKind()
		key := client.ObjectKey(l.resource)
		prefix = fmt.Sprintf("%s %s/%s | %s |", prefix, sprint(gvk.GroupVersion()), sprint(gvk.Kind), client.ColouredName(key, color))
	}
	a = append(a, prefix)
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
