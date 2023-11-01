package operations

import (
	"bytes"
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	fakeLogger "github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/stretchr/testify/assert"
)

type FakeCommandOutput struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (c *FakeCommandOutput) Out() string {
	return c.stdout.String()
}

func (c *FakeCommandOutput) Err() string {
	return c.stderr.String()
}

func TestExec(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name      string
		exec      v1alpha1.Exec
		log       bool
		namespace string
		expected  []string
		err       error
	}{
		{
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
		},
		{
			name: "Test with Script",
			exec: v1alpha1.Exec{
				Script: &v1alpha1.Script{
					Content: "echo hello",
				},
			},
			log:       true,
			namespace: "test-namespace",
			expected:  []string{"SCRIPT: [RUNNING...]", "STDOUT: [LOGS...\nhello]", "SCRIPT: [DONE]"},
		},
		{
			name:      "Test with nil Command and Script",
			exec:      v1alpha1.Exec{},
			log:       true,
			namespace: "test-namespace",
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &fakeLogger.MockLogger{}
			Exec(ctx, logger, tt.exec, tt.log, tt.namespace)
			assert.ElementsMatch(t, tt.expected, logger.Logs)
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
	}{
		{
			name: "Test with valid Command",
			command: v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"hello"},
			},
			log:       true,
			namespace: "test-namespace",
			wantErr:   false,
		},
		{
			name: "Test with invalid Command",
			command: v1alpha1.Command{
				Entrypoint: "invalidCmd",
			},
			log:       true,
			namespace: "test-namespace",
			wantErr:   true,
		},
		{
			name: "Test without logging",
			command: v1alpha1.Command{
				Entrypoint: "echo",
				Args:       []string{"silent"},
			},
			log:       false,
			namespace: "test-namespace",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := command(ctx, &fakeLogger.MockLogger{}, tt.command, tt.log, tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("command() error = %v, wantErr %v", err, tt.wantErr)
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
	}{
		{
			name: "Test with valid Script",
			script: v1alpha1.Script{
				Content: "echo hello",
			},
			log:       true,
			namespace: "test-namespace",
			wantErr:   false,
		},
		{
			name: "Test with invalid Script",
			script: v1alpha1.Script{
				Content: "invalidScriptCommand",
			},
			log:       true,
			namespace: "test-namespace",
			wantErr:   true,
		},
		{
			name: "Test script without logging",
			script: v1alpha1.Script{
				Content: "echo silent",
			},
			log:       false,
			namespace: "test-namespace",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := script(ctx, &fakeLogger.MockLogger{}, tt.script, tt.log, tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("script() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
