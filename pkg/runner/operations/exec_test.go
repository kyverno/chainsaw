package operations

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/stretchr/testify/assert"
)

func Test_operationExec(t *testing.T) {
	tests := []struct {
		name      string
		exec      v1alpha1.Exec
		log       bool
		namespace string
		expected  []string
		wantErr   bool
	}{{
		name: "Command execution",
		exec: v1alpha1.Exec{
			Command: &v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"hello"},
			},
		},
		log:       true,
		namespace: "test-namespace",
		expected: []string{
			"CMD   : [/usr/bin/echo hello RUNNING...]",
			"STDOUT: [LOGS...\nhello]",
			"CMD   : [DONE]",
		},
	}, {
		name: "Test with Script",
		exec: v1alpha1.Exec{
			Script: &v1alpha1.Script{
				Content: "echo hello",
			},
		},
		log:       true,
		namespace: "test-namespace",
		expected:  []string{"SCRIPT: [RUNNING...]", "STDOUT: [LOGS...\nhello]", "SCRIPT: [DONE]"},
	}, {
		name:      "Test with nil Command and Script",
		exec:      v1alpha1.Exec{},
		log:       true,
		namespace: "test-namespace",
		expected:  nil,
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &MockLogger{}
			ctx := logging.IntoContext(context.TODO(), logger)
			err := operationExec(ctx, tt.exec, tt.log, tt.namespace)
			assert.ElementsMatch(t, tt.expected, logger.Logs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommand(t *testing.T) {
	tests := []struct {
		name      string
		command   v1alpha1.Command
		log       bool
		namespace string
		wantErr   bool
	}{{
		name: "Test with valid Command",
		command: v1alpha1.Command{
			Entrypoint: "echo",
			Args:       []string{"hello"},
		},
		log:       true,
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test with invalid Command",
		command: v1alpha1.Command{
			Entrypoint: "invalidCmd",
		},
		log:       true,
		namespace: "test-namespace",
		wantErr:   true,
	}, {
		name: "Test without logging",
		command: v1alpha1.Command{
			Entrypoint: "echo",
			Args:       []string{"silent"},
		},
		log:       false,
		namespace: "test-namespace",
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logging.IntoContext(context.TODO(), &MockLogger{})
			err := command(ctx, tt.command, tt.log, tt.namespace)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScript(t *testing.T) {
	tests := []struct {
		name      string
		script    v1alpha1.Script
		log       bool
		namespace string
		wantErr   bool
	}{{
		name: "Test with valid Script",
		script: v1alpha1.Script{
			Content: "echo hello",
		},
		log:       true,
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test with invalid Script",
		script: v1alpha1.Script{
			Content: "invalidScriptCommand",
		},
		log:       true,
		namespace: "test-namespace",
		wantErr:   true,
	}, {
		name: "Test script without logging",
		script: v1alpha1.Script{
			Content: "echo silent",
		},
		log:       false,
		namespace: "test-namespace",
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logging.IntoContext(context.TODO(), &MockLogger{})
			err := script(ctx, tt.script, tt.log, tt.namespace)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_expand(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		in   []string
		want []string
	}{{
		name: "nil",
		env:  nil,
		in:   []string{"echo", "$NAMESPACE"},
		want: []string{"echo", "$NAMESPACE"},
	}, {
		name: "empty",
		env:  map[string]string{},
		in:   []string{"echo", "$NAMESPACE"},
		want: []string{"echo", "$NAMESPACE"},
	}, {
		name: "expand",
		env:  map[string]string{"NAMESPACE": "foo"},
		in:   []string{"echo", "$NAMESPACE"},
		want: []string{"echo", "foo"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expand(tt.env, tt.in...)
			assert.Equal(t, tt.want, got)
		})
	}
}
