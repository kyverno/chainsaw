package logging

import (
	"context"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestLoggerFunc_Log(t *testing.T) {
	tests := []struct {
		name      string
		test      string
		scenario  string
		step      string
		operation Operation
		status    Status
		obj       client.Object
		color     *color.Color
		args      []fmt.Stringer
	}{{
		test:      "test",
		step:      "step",
		operation: Create,
		status:    BeginStatus,
		obj:       &unstructured.Unstructured{},
		color:     color.New(color.FgBlue),
		args:      []fmt.Stringer{Section("arg1"), Section("arg2")},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			var f SinkFunc = func(test, scenario, step string, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
				assert.Equal(t, tt.test, test)
				assert.Equal(t, tt.scenario, scenario)
				assert.Equal(t, tt.step, step)
				assert.Equal(t, tt.operation, operation)
				assert.Equal(t, tt.status, status)
				assert.Equal(t, tt.obj, obj)
				assert.Equal(t, tt.color, color)
				assert.Equal(t, tt.args, args)
				called = true
			}
			ctx := context.Background()
			ctx = WithSink(ctx, f)
			logger := NewLogger(tt.test, tt.scenario, tt.step)
			logger.Log(ctx, tt.operation, tt.status, tt.obj, tt.color, tt.args...)
			assert.True(t, called)
		})
	}
}
