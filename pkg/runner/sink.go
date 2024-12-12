package runner

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/testing"
	"k8s.io/utils/clock"
)

const eraser = "\b\b\b\b\b\b\b\b\b\b\b\b"

func newSink(clock clock.PassiveClock, t testing.TTest) logging.SinkFunc {
	return func(test string, step string, operation logging.Operation, status logging.Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
		sprint := fmt.Sprint
		opLen := 9
		stLen := 5
		if color != nil {
			sprint = color.Sprint
			opLen += 14
			stLen += 14
		}
		a := make([]any, 0, len(args)+2)
		prefix := fmt.Sprintf("%s| %s | %s | %s | %-*s | %-*s |", eraser, clock.Now().Format("15:04:05"), sprint(test), sprint(step), opLen, sprint(operation), stLen, sprint(status))
		if obj != nil {
			gvk := obj.GetObjectKind().GroupVersionKind()
			key := client.Key(obj)
			prefix = fmt.Sprintf("%s %s/%s @ %s", prefix, gvk.GroupVersion(), gvk.Kind, client.Name(key))
		}
		a = append(a, prefix)
		for _, arg := range args {
			a = append(a, "\n")
			a = append(a, arg)
		}
		t.Log(fmt.Sprint(a...))
	}
}
