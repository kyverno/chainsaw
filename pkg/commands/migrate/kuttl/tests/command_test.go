package tests

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/kyverno/chainsaw/pkg/commands/root"
	"github.com/stretchr/testify/assert"
)

const kuttlOutPath = "../../../../../testdata/kuttl"

func Test_Execute(t *testing.T) {
	basePath := "../../../../../testdata/commands/migrate/kuttl/tests"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{{
		name: "help",
		args: []string{
			"tests",
			"--help",
		},
		out:     filepath.Join(basePath, "help.txt"),
		wantErr: false,
	}, {
		name: "migrate",
		args: []string{
			"tests",
			kuttlOutPath,
		},
		out:     filepath.Join(basePath, "out.txt"),
		wantErr: false,
	}, {
		name: "migrate save",
		args: []string{
			"tests",
			kuttlOutPath,
			"--save",
		},
		out:     filepath.Join(basePath, "out-save.txt"),
		wantErr: false,
	}, {
		name: "unknow flag",
		args: []string{
			"tests",
			"--foo",
		},
		wantErr: true,
	}, {
		name: "bad folder",
		args: []string{
			"tests",
			petname.Name(),
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

			if slices.Contains(tt.args, "--save") {
				defer os.Remove(filepath.Join(kuttlOutPath, "chainsaw-test.yaml"))
			}
			cmd.SetErr(os.Stderr)

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
