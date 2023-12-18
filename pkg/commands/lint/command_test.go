package lint

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/commands/root"
	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	basePath := "../../../testdata/commands/lint"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{
		{
			name: "Lint Test JSON File",
			args: []string{
				"lint",
				"test",
				"--file",
				filepath.Join(basePath, "test", "test.json"),
			},
			out:     filepath.Join(basePath, "test", "pass.txt"),
			wantErr: false,
		},
		{
			name: "Lint Test YAML File",
			args: []string{
				"lint",
				"test",
				"--file",
				filepath.Join(basePath, "test", "test.yaml"),
			},
			out:     filepath.Join(basePath, "test", "pass.txt"),
			wantErr: false,
		},
		{
			name: "Lint Configuration JSON File",
			args: []string{
				"lint",
				"configuration",
				"--file",
				filepath.Join(basePath, "configuration", "configuration.json"),
			},
			out:     filepath.Join(basePath, "configuration", "pass.txt"),
			wantErr: false,
		},
		{
			name: "Lint Configuration YAML File",
			args: []string{
				"lint",
				"configuration",
				"--file",
				filepath.Join(basePath, "configuration", "configuration.yaml"),
			},
			out:     filepath.Join(basePath, "configuration", "pass.txt"),
			wantErr: false,
		},
		{
			name: "Lint Error Test JSON File",
			args: []string{
				"lint",
				"test",
				"--file",
				filepath.Join(basePath, "test", "wrong-test.json"),
			},
			wantErr: true,
		},
		{
			name: "Lint Error Test YAML File",
			args: []string{
				"lint",
				"test",
				"--file",
				filepath.Join(basePath, "test", "wrong-test.yaml"),
			},
			wantErr: true,
		},
		{
			name: "Lint Error Configuration JSON File",
			args: []string{
				"lint",
				"configuration",
				"--file",
				filepath.Join(basePath, "configuration", "wrong-configuration.json"),
			},
			wantErr: true,
		},
		{
			name: "Lint Error Configuration YAML File",
			args: []string{
				"lint",
				"configuration",
				"--file",
				filepath.Join(basePath, "configuration", "wrong-configuration.yaml"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := root.Command()
			cmd.AddCommand(Command())
			assert.NotNil(t, cmd)
			cmd.SetArgs(tt.args)
			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			actual, err := io.ReadAll(out)
			assert.NoError(t, err)
			if tt.out != "" {
				expected, err := os.ReadFile(tt.out)
				assert.NoError(t, err)
				assert.Equal(t, string(expected), string(actual))
			}
		})
	}
}
