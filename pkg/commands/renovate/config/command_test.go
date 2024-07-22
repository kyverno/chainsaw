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
	basePath := "../../../../testdata/commands/renovate/config"
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
		name: "v1alpha1/default",
		args: []string{
			"config",
			"../../../../testdata/config/v1alpha1/default.yaml",
		},
		out:     filepath.Join(basePath, "v1alpha1-default.txt"),
		wantErr: false,
	}, {
		name: "v1alpha1/custom",
		args: []string{
			"config",
			"../../../../testdata/config/v1alpha1/custom-config.yaml",
		},
		out:     filepath.Join(basePath, "v1alpha1-custom.txt"),
		wantErr: false,
		// }, {
		// 	name: "renovate save",
		// 	args: []string{
		// 		"config",
		// 		"../../../../testdata/config/v1alpha1/default.yaml",
		// 		"--save",
		// 	},
		// 	out:     filepath.Join(basePath, "out-save.txt"),
		// 	wantErr: false,
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
			"../../../../testdata/config/unknown.yaml",
		},
		wantErr: true,
	}, {
		name: "multiple file",
		args: []string{
			"config",
			"../../../../testdata/config/multiple.yaml",
		},
		wantErr: true,
	}, {
		name: "invalid config",
		args: []string{
			"config",
			"../../../../testdata/config/v1alpha1/bad-catch.yaml",
		},
		wantErr: true,
	}, {
		name: "configmap",
		args: []string{
			"config",
			"../../../../testdata/config/configmap.yaml",
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
