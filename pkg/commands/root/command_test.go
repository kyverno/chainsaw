package root

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	basePath := "../../../testdata/commands/root"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{{
		name: "help",
		args: []string{
			"--help",
		},
		out:     filepath.Join(basePath, "help.txt"),
		wantErr: false,
	}, {
		name:    "chainsaw",
		out:     filepath.Join(basePath, "help.txt"),
		wantErr: false,
	}, {
		name:    "unknow flag",
		args:    []string{"--foo"},
		wantErr: true,
	}}
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
