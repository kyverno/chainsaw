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

func Test_NewSink(t *testing.T) {
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
			got := NewSink(tclock.NewFakePassiveClock(time.Time{}), log)
			if tt.color != nil {
				tt.color.EnableColor()
			}
			got.Log(tt.test, tt.step, tt.operation, tt.status, tt.obj, tt.color, tt.args...)
			assert.Equal(t, tt.want, out)
		})
	}
}

func Test_NewFilteredSink(t *testing.T) {
	tests := []struct {
		name       string
		test       string
		step       string
		operation  logging.Operation
		status     logging.Status
		obj        client.Object
		color      *color.Color
		args       []fmt.Stringer
		noWarnings bool
		wantOutput bool
	}{
		{
			name:       "passes non-warning status",
			test:       "foo",
			step:       "bar",
			operation:  logging.Apply,
			status:     logging.OkStatus,
			noWarnings: true,
			wantOutput: true,
		},
		{
			name:       "filters warning status",
			test:       "foo",
			step:       "bar",
			operation:  logging.Apply,
			status:     logging.WarnStatus,
			noWarnings: true,
			wantOutput: false,
		},
		{
			name:       "passes warning status when noWarnings is false",
			test:       "foo",
			step:       "bar",
			operation:  logging.Apply,
			status:     logging.WarnStatus,
			noWarnings: false,
			wantOutput: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			log := func(args ...any) {
				called = true
			}
			got := NewFilteredSink(tclock.NewFakePassiveClock(time.Time{}), log, tt.noWarnings)
			got.Log(tt.test, tt.step, tt.operation, tt.status, tt.obj, tt.color, tt.args...)
			assert.Equal(t, tt.wantOutput, called)
		})
	}
}
