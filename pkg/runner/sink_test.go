package runner

import (
	"fmt"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/utils/kube"
	"github.com/stretchr/testify/assert"
	tclock "k8s.io/utils/clock/testing"
	"k8s.io/utils/ptr"
)

func Test_newSink(t *testing.T) {
	tests := []struct {
		name      string
		test      string
		step      string
		operation logging.Operation
		status    logging.Status
		obj       client.Object
		color     *color.Color
		args      []fmt.Stringer
		want      []any
	}{{
		name:      "simple",
		test:      "foo",
		step:      "bar",
		operation: logging.Apply,
		status:    logging.OkStatus,
		obj:       nil,
		color:     nil,
		args:      nil,
		want:      []any{"\b\b\b\b\b\b\b\b\b\b\b\b| 00:00:00 | foo | bar | APPLY     | OK    |"},
	}, {
		name:      "with color",
		test:      "foo",
		step:      "bar",
		operation: logging.Apply,
		status:    logging.OkStatus,
		obj:       nil,
		color:     color.New(color.FgGreen),
		args:      nil,
		want:      []any{"\b\b\b\b\b\b\b\b\b\b\b\b| 00:00:00 | \x1b[32mfoo\x1b[0m | \x1b[32mbar\x1b[0m | \x1b[32mAPPLY\x1b[0m          | \x1b[32mOK\x1b[0m         |"},
	}, {
		name:      "with object",
		test:      "foo",
		step:      "bar",
		operation: logging.Apply,
		status:    logging.OkStatus,
		obj:       ptr.To(kube.Namespace("dog")),
		color:     nil,
		args:      nil,
		want:      []any{"\b\b\b\b\b\b\b\b\b\b\b\b| 00:00:00 | foo | bar | APPLY     | OK    | v1/Namespace @ dog"},
	}, {
		name:      "with color and object",
		test:      "foo",
		step:      "bar",
		operation: logging.Apply,
		status:    logging.OkStatus,
		obj:       ptr.To(kube.Namespace("dog")),
		color:     color.New(color.FgGreen),
		args:      nil,
		want:      []any{"\b\b\b\b\b\b\b\b\b\b\b\b| 00:00:00 | \x1b[32mfoo\x1b[0m | \x1b[32mbar\x1b[0m | \x1b[32mAPPLY\x1b[0m          | \x1b[32mOK\x1b[0m         | v1/Namespace @ dog"},
	}, {
		name:      "with args",
		test:      "foo",
		step:      "bar",
		operation: logging.Apply,
		status:    logging.OkStatus,
		obj:       nil,
		color:     nil,
		args:      []fmt.Stringer{logging.Section("foo", "bar")},
		want:      []any{"\b\b\b\b\b\b\b\b\b\b\b\b| 00:00:00 | foo | bar | APPLY     | OK    |\n=== FOO\nbar"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out []any
			log := func(args ...any) {
				out = args
			}
			got := newSink(tclock.NewFakePassiveClock(time.Time{}), false, log)
			if tt.color != nil {
				tt.color.EnableColor()
			}
			got.Log(tt.test, tt.step, tt.operation, tt.status, tt.obj, tt.color, tt.args...)
			assert.Equal(t, tt.want, out)
		})
	}
}
