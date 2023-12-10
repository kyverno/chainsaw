package docs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	basePath := "../../../testdata/commands/docs"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{
		{
			name:    "help",
			args:    []string{"docs", "--help"},
			wantErr: false,
			out:     filepath.Join(basePath, "help.txt"),
		},
		{
			name:    "generate documentation",
			args:    []string{"docs", "--test-file", "example-test.yaml", "--test-dir", basePath},
			wantErr: false,
			out:     filepath.Join(basePath, "invalid-output.txt"),
		},
		{
			name:    "unknown flag",
			args:    []string{"docs", "--foo"},
			wantErr: true,
		},
		{
			name:    "unknown arg",
			args:    []string{"docs", "foo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command()
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
