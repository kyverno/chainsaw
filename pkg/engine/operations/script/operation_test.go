package script

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func Test_operationScript(t *testing.T) {
	tests := []struct {
		name       string
		script     v1alpha1.Script
		basePath   string
		namespace  string
		wantErr    bool
		wantErrMsg string
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
		name: "Test with absolute workdir",
		script: v1alpha1.Script{
			Content:   "cat operation.go",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
			WorkDir:   ptr.To("/bar"),
		},
		basePath:   "..",
		namespace:  "test-namespace",
		wantErr:    true,
		wantErrMsg: "/bar: no such file or directory",
	}, {
		name: "Test with relative workdir",
		script: v1alpha1.Script{
			Content:   "cat operation.go",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
			WorkDir:   ptr.To("./foo"),
		},
		basePath:   "..",
		namespace:  "test-namespace",
		wantErr:    true,
		wantErrMsg: "../foo: no such file or directory",
	}, {
		name: "with check",
		script: v1alpha1.Script{
			Content:   "foo",
			ActionEnv: v1alpha1.ActionEnv{SkipLogOutput: true},
			ActionCheck: v1alpha1.ActionCheck{
				Check: ptr.To(v1alpha1.NewCheck(
					map[string]any{
						"($error != null)": true,
					},
				)),
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
				Check: ptr.To(v1alpha1.NewCheck(
					map[string]any{
						"(foo('bar'))": true,
					},
				)),
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
				Check: ptr.To(v1alpha1.NewCheck(
					map[string]any{
						"(foo('bar'))": true,
					},
				)),
			},
		},
		basePath:  "..",
		namespace: "test-namespace",
		wantErr:   true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := logging.WithLogger(context.TODO(), &mocks.Logger{})
			operation := New(
				apis.DefaultCompilers,
				tt.script,
				tt.basePath,
				tt.namespace,
				nil,
			)
			_, err := operation.Exec(ctx, nil)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil && tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
