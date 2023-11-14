package operations

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
)

func Test_operationCommand(t *testing.T) {
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
			ctx := logging.IntoContext(context.TODO(), &tlogging.FakeLogger{})
			commandOp := &CommandOperation{
				command:       tt.command,
				skipLogOutput: tt.log,
				namespace:     tt.namespace,
			}
			err := execOperation(ctx, commandOp)
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
