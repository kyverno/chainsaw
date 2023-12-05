package config

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
	basePath := "../../../../../testdata/commands/migrate/kuttl/config"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{{
		name: "help",
		args: []string{
			"config",
			"--help",
		},
		out:     filepath.Join(basePath, "help.txt"),
		wantErr: false,
	}, {
		name: "migrate",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/kuttl-test.yaml",
		},
		out:     filepath.Join(basePath, "out.txt"),
		wantErr: false,
	}, {
		name: "migrate save",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/kuttl-test.yaml",
			"--save",
		},
		out:     filepath.Join(basePath, "out-save.txt"),
		wantErr: false,
	}, {
		name: "unknow flag",
		args: []string{
			"config",
			"--foo",
		},
		wantErr: true,
	}, {
		name: "unknown file",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/unknown.yaml",
		},
		wantErr: true,
	}, {
		name: "multiple file",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/multiple-config.yaml",
		},
		wantErr: true,
	}, {
		name: "not a config",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/02-step.yaml",
		},
		wantErr: true,
	}, {
		name: "invalid config",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/invalid-config.yaml",
		},
		wantErr: true,
	}, {
		name: "configmap",
		args: []string{
			"config",
			"../../../../../testdata/kuttl/configmap.yaml",
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
