package logging

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	tclock "k8s.io/utils/clock/testing"
)

func TestNewLogger(t *testing.T) {
	fakeClock := tclock.NewFakePassiveClock(time.Now())
	testName := "testName"
	stepName := "stepName"
	logger, ok := NewLogger(t, fakeClock, testName, stepName).(*logger)

	assert.True(t, ok, "Type assertion for *logger failed")

	assert.Equal(t, t, logger.t)
	assert.Equal(t, fakeClock, logger.clock)
	assert.Equal(t, testName, logger.test)
	assert.Equal(t, stepName, logger.step)
	assert.Nil(t, logger.resource)
}

type s string

func (v s) String() string { return string(v) }

func Test_logger_Log(t *testing.T) {
	fakeClock := tclock.NewFakePassiveClock(time.Now())
	mockT := &tlogging.FakeTLogger{}
	fakeLogger := NewLogger(mockT, fakeClock, "testName", "stepName").(*logger)
	disabled := color.New(color.FgBlue)
	disabled.DisableColor()
	enabled := color.New(color.FgBlue)
	enabled.EnableColor()
	testCases := []struct {
		name           string
		resource       client.Object
		operation      string
		status         string
		color          *color.Color
		args           []fmt.Stringer
		expectContains []string
	}{
		{
			name:      "without resource",
			resource:  nil,
			operation: "OPERATION",
			status:    "STATUS",
			args:      []fmt.Stringer{s("arg1"), s("arg2")},
			expectContains: []string{
				"testName", "stepName", "OPERATION", "arg1", "arg2",
			},
		},
		{
			name:      "with color",
			resource:  nil,
			operation: "OPERATION",
			status:    "STATUS",
			color:     enabled,
			args:      []fmt.Stringer{s("arg1"), s("arg2")},
			expectContains: []string{
				"testName", "stepName", "OPERATION", "arg1", "arg2",
			},
		},
		{
			name: "with resource",
			resource: func() client.Object {
				var r unstructured.Unstructured
				r.SetName("testResource")
				r.SetNamespace("default")
				r.SetAPIVersion("testGroup/v1")
				r.SetKind("testKind")
				return &r
			}(),
			operation: "OPERATION",
			status:    "STATUS",
			args:      []fmt.Stringer{s("arg1"), s("arg2")},
			expectContains: []string{
				"testName", "stepName", "OPERATION", "default/testResource", "testGroup/v1/testKind", "arg1", "arg2",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resource != nil {
				fakeLogger = fakeLogger.WithResource(tt.resource).(*logger)
			}
			fakeLogger.Log(Operation(tt.operation), Status(tt.status), tt.color, tt.args...)
			for _, exp := range tt.expectContains {
				found := false
				for _, msg := range mockT.Messages {
					if strings.Contains(msg, exp) {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected to find '%s' in logs, but didn't. Logs: %v", exp, mockT.Messages)
			}
			mockT.Messages = []string{}
		})
	}
}

func Test_logger_WithResource(t *testing.T) {
	testCases := []struct {
		name      string
		resource  client.Object
		expectNil bool
	}{{
		name: "Valid Resource",
		resource: func() client.Object {
			var r unstructured.Unstructured
			r.SetName("testResource")
			return &r
		}(),
		expectNil: false,
	}, {
		name:      "Nil Resource",
		resource:  nil,
		expectNil: true,
	}}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fakeClock := tclock.NewFakePassiveClock(time.Now())
			fakeLogger := logger{
				t:     t,
				clock: fakeClock,
				test:  "testName",
				step:  "stepName",
			}

			newLogger := fakeLogger.WithResource(tt.resource).(*logger)

			if tt.expectNil {
				assert.Nil(t, newLogger.resource, "Expected resource to be nil in the logger")
			} else {
				assert.NotNil(t, newLogger.resource, "Expected resource to not be nil in the logger")
				assert.Equal(t, tt.resource, newLogger.resource, "Expected correct resource to be set in the logger")
			}

			assert.Equal(t, fakeLogger.t, newLogger.t, "Expected testing.T to remain the same")
			assert.Equal(t, fakeLogger.clock, newLogger.clock, "Expected clock to remain the same")
			assert.Equal(t, fakeLogger.test, newLogger.test, "Expected test name to remain the same")
			assert.Equal(t, fakeLogger.step, newLogger.step, "Expected step name to remain the same")
		})
	}
}
