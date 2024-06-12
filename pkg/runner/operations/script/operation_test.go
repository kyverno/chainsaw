package script

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
		basePath  string
		namespace string
		wantErr   bool
	}{{
		name: "Test with valid Script",
		script: v1alpha1.Script{
			Content:   "echo hello",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: false},
		},
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test with invalid Script",
		script: v1alpha1.Script{
			Content:   "invalidScriptCommand",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: false},
		},
		namespace: "test-namespace",
		wantErr:   true,
	}, {
		name: "Test script without logging",
		script: v1alpha1.Script{
			Content:   "echo silent",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
		},
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test base path",
		script: v1alpha1.Script{
			Content:   "cat operation.go",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
		},
		basePath:  "..",
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "with check",
		script: v1alpha1.Script{
			Content:   "foo",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
			ActionCheck: v1alpha1.ActionCheck{
				Check: &v1alpha1.Check{
					Value: map[string]any{
						"($error != null)": true,
					},
				},
			},
		},
		basePath:  "..",
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "with bad check",
		script: v1alpha1.Script{
			Content:   "foo",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
			ActionCheck: v1alpha1.ActionCheck{
				Check: &v1alpha1.Check{
					Value: map[string]any{
						"(foo('bar'))": true,
					},
				},
			},
		},
		basePath:  "..",
		namespace: "test-namespace",
		wantErr:   true,
	}, {
		name: "with bad check",
		script: v1alpha1.Script{
			Content:   "cat operation.go",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
			ActionCheck: v1alpha1.ActionCheck{
				Check: &v1alpha1.Check{
					Value: map[string]any{
						"(foo('bar'))": true,
					},
				},
			},
		},
		basePath:  "..",
		namespace: "test-namespace",
		wantErr:   true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logging.IntoContext(context.TODO(), &tlogging.FakeLogger{})
			operation := New(
				tt.script,
				tt.basePath,
				tt.namespace,
				nil,
			)
			_, err := operation.Exec(ctx, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
