package command

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func Test_operationCommand(t *testing.T) {
	tests := []struct {
		name       string
		command    v1alpha1.Command
		basePath   string
		namespace  string
		wantErr    bool
		wantErrMsg string
	}{{
		name: "Test with valid Command",
		command: v1alpha1.Command{
			Entrypoint: "echo",
			Args:       []string{"hello"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: false},
		},
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test with invalid Command",
		command: v1alpha1.Command{
			Entrypoint: "invalidCmd",
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: false},
		},
		namespace: "test-namespace",
		wantErr:   true,
	}, {
		name: "Test without logging",
		command: v1alpha1.Command{
			Entrypoint: "echo",
			Args:       []string{"silent"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
		},
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test base path",
		command: v1alpha1.Command{
			Entrypoint: "cat",
			Args:       []string{"operation.go"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
		},
		basePath:  "..",
		namespace: "test-namespace",
		wantErr:   false,
	}, {
		name: "Test with absolute workdir",
		command: v1alpha1.Command{
			Entrypoint: "cat",
			Args:       []string{"operation.go"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
			WorkDir:    ptr.To("/bar"),
		},
		basePath:   "..",
		namespace:  "test-namespace",
		wantErr:    true,
		wantErrMsg: "/bar: no such file or directory",
	}, {
		name: "Test with relative workdir",
		command: v1alpha1.Command{
			Entrypoint: "cat",
			Args:       []string{"operation.go"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
			WorkDir:    ptr.To("./foo"),
		},
		basePath:   "..",
		namespace:  "test-namespace",
		wantErr:    true,
		wantErrMsg: "../foo: no such file or directory",
	}, {
		name: "with check",
		command: v1alpha1.Command{
			Entrypoint: "foo",
			Args:       []string{"operation.go"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
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
		command: v1alpha1.Command{
			Entrypoint: "foo",
			Args:       []string{"operation.go"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
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
		command: v1alpha1.Command{
			Entrypoint: "cat",
			Args:       []string{"operation.go"},
			ActionEnv:  v1alpha1.ActionEnv{SkipLogOutput: true},
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
			ctx := logging.IntoContext(context.TODO(), &tlogging.FakeLogger{})
			operation := New(
				tt.command,
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
