package operations

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
)

func Test_operationScript(t *testing.T) {
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
			ctx := logging.IntoContext(context.TODO(), &tlogging.FakeLogger{})
			scriptOp := &ScriptOperation{
				script:        tt.script,
				skipLogOutput: tt.log,
				namespace:     tt.namespace,
			}
			err := execOperation(ctx, scriptOp)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
