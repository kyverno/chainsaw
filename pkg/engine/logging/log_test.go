package logging

import (
	"context"
	"fmt"
	"testing"
	"time"

	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/kyverno/pkg/ext/output/color"
	tclock "k8s.io/utils/clock/testing"
)

func TestLog(t *testing.T) {
	fakeClock := tclock.NewFakePassiveClock(time.Now())
	mockT := &tlogging.FakeTLogger{}
	fakeLogger := NewLogger(mockT, fakeClock, "testName", "stepName").(*logger)
	tests := []struct {
		name      string
		ctx       context.Context //nolint:containedctx
		operation string
		status    string
		color     *color.Color
		args      []fmt.Stringer
	}{{
		name:      "background",
		ctx:       context.Background(),
		operation: "foo",
		status:    "bar",
		color:     nil,
		args:      nil,
	}, {
		name:      "nil",
		ctx:       nil,
		operation: "foo",
		status:    "bar",
		color:     nil,
		args:      nil,
	}, {
		name:      "with logger",
		ctx:       IntoContext(context.Background(), fakeLogger),
		operation: "foo",
		status:    "bar",
		color:     nil,
		args:      nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.ctx, Operation(tt.operation), Status(tt.status), tt.color, tt.args...)
		})
	}
}
