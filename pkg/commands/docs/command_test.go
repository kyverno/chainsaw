package docs

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
	basePath := "../../../testdata/commands/docs"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{{
		name: "help",
		args: []string{
			"docs",
			"--help",
		},
		out:     filepath.Join(basePath, "help.txt"),
		wantErr: false,
	}, {
		name: "invalid output",
		args: []string{
			"docs",
		},
		out:     filepath.Join(basePath, "invalid-output.txt"),
		wantErr: true,
	}, {
		name: "docs",
		args: []string{
			"docs",
			"--output",
			"/tmp",
		},
		wantErr: false,
	}, {
		name: "website docs",
		args: []string{
			"docs",
			"--website",
			"--output",
			"/tmp",
		},
		wantErr: false,
	}, {
		name: "unknow flag",
		args: []string{
			"docs",
			"--foo",
		},
		wantErr: true,
	}, {
		name: "unknow arg",
		args: []string{
			"docs",
			"foo",
		},
		wantErr: true,
	}}
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
