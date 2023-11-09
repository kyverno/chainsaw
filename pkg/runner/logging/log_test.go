package logging

import (
	"context"
	"testing"
	"time"

	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/kyverno/kyverno/ext/output/color"
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
		color     *color.Color
		args      []interface{}
	}{{
		name:      "background",
		ctx:       context.Background(),
		operation: "foo",
		color:     nil,
		args:      nil,
	}, {
		name:      "nil",
		ctx:       nil,
		operation: "foo",
		color:     nil,
		args:      nil,
	}, {
		name:      "with logger",
		ctx:       IntoContext(context.Background(), fakeLogger),
		operation: "foo",
		color:     nil,
		args:      nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.ctx, tt.operation, tt.color, tt.args...)
		})
	}
}
