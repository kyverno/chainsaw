package test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainsawCommand(t *testing.T) {
	basePath := "../../../testdata/commands/test"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
		err     string
	}{
		{
			name:    "default test",
			args:    []string{},
			wantErr: false,
			out:     filepath.Join(basePath, "default.txt"),
		},
		{
			name: "with timeout",
			args: []string{
				"--timeout",
				"10s",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "without_config.txt"),
		},
		{
			name: "invalid timeout",
			args: []string{
				"--timeout",
				"invalid",
			},
			wantErr: true,
		},
		{
			name: "test dirs specified",
			args: []string{
				"--testDirs",
				"dir1,dir2,dir3",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "without_config.txt"),
		},
		{
			name: "nonexistent config file",
			args: []string{
				"--config",
				"nonexistent.yaml",
			},
			wantErr: true,
		},
		{
			name: "suppress logs",
			args: []string{
				"--suppress",
				"warning,error",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "without_config.txt"),
		},
		{
			name: "skip test with regex",
			args: []string{
				"--skipTestRegex",
				"test[1-3]",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "without_config.txt"),
		},
		{
			name: "valid config",
			args: []string{
				"--config",
				filepath.Join(basePath, "config/empty_config.yaml"),
			},
			wantErr: false,
			out:     filepath.Join(basePath, "valid_config.txt"),
		},
		{
			name: "nonexistent config",
			args: []string{
				"--config",
				filepath.Join(basePath, "config/nonexistent_config.yaml"),
			},
			wantErr: true,
		},
		{
			name: "misformatted config",
			args: []string{
				"--config",
				filepath.Join(basePath, "config/wrong_format_config.yaml"),
			},
			wantErr: true,
			out:     filepath.Join(basePath, "wrong_format_config.txt"),
		},
		{
			name: "wrong kind in config",
			args: []string{
				"--config",
				filepath.Join(basePath, "config/wrong_kind_config.yaml"),
			},
			wantErr: true,
			out:     filepath.Join(basePath, "wrong_kind_config.txt"),
			err:     filepath.Join(basePath, "wrong_kind_config_err.txt"),
		},
		{
			name: "config with all fields",
			args: []string{
				"--config",
				filepath.Join(basePath, "config/config_all_fields.yaml"),
				"--timeout",
				"10s",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "config_all_fields.txt"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command()
			assert.NotNil(t, cmd)
			cmd.SetArgs(tt.args)
			stdout := bytes.NewBufferString("")
			cmd.SetOut(stdout)
			stderr := bytes.NewBufferString("")
			cmd.SetErr(stderr)
			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			actualOut, err := io.ReadAll(stdout)
			assert.NoError(t, err)
			actualErr, err := io.ReadAll(stderr)
			assert.NoError(t, err)
			if tt.out != "" {
				expected, err := os.ReadFile(tt.out)
				assert.NoError(t, err)
				assert.Equal(t, string(expected), string(actualOut))
			}
			if tt.err != "" {
				expected, err := os.ReadFile(tt.err)
				assert.NoError(t, err)
				assert.Equal(t, string(expected), string(actualErr))
			}
		})
	}
}
