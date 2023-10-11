package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	}{
		{
			name:    "default test",
			args:    []string{},
			wantErr: false,
			out:     basePath + "/default.txt",
		},
		{
			name:    "with duration",
			args:    []string{"--duration=10s"},
			wantErr: false,
			out:     basePath + "/without_config.txt",
		},
		{
			name:    "invalid duration",
			args:    []string{"--duration=invalid"},
			wantErr: true,
		},
		{
			name:    "test dirs specified",
			args:    []string{"--testDirs=dir1,dir2,dir3"},
			wantErr: false,
			out:     basePath + "/without_config.txt",
		},
		{
			name:    "nonexistent config file",
			args:    []string{"--config=nonexistent.yaml"},
			wantErr: true,
		},
		{
			name:    "suppress logs",
			args:    []string{"--suppress=warning,error"},
			wantErr: false,
			out:     basePath + "/without_config.txt",
		},
		{
			name:    "skip test with regex",
			args:    []string{"--skipTestRegex=test[1-3]"},
			wantErr: false,
			out:     basePath + "/without_config.txt",
		},
		{
			name:    "valid config",
			args:    []string{fmt.Sprintf("--config=%s/config/empty_config.yaml", basePath)},
			wantErr: false,
			out:     basePath + "/valid_config.txt",
		},
		{
			name:    "nonexistent config",
			args:    []string{fmt.Sprintf("--config=%s/config/nonexistent_config.yaml", basePath)},
			wantErr: true,
		},

		// fix

		// {
		// 	name:    "misformatted config",
		// 	args:    []string{fmt.Sprintf("--config=%s/config/wrong_format_config.yaml", basePath)},
		// 	wantErr: true,
		// 	out:     basePath + "/wrong_format_config.txt",
		// },
		{
			name:    "wrong kind in config",
			args:    []string{fmt.Sprintf("--config=%s/config/wrong_kind_config.yaml", basePath)},
			wantErr: true,
			out:     basePath + "/wrong_kind_config.txt",
		},
		// {
		// 	name:    "valid config with duration",
		// 	args:    []string{"--config=config.yaml", "--duration=10s"},
		// 	wantErr: false,
		// 	out:     basePath + "valid_config_with_duration.txt",
		// },
		// {
		// 	name:    "invalid config with duration",
		// 	args:    []string{"--config=wrong_config.yaml", "--duration=10s"},
		// 	wantErr: true,
		// 	out:     basePath + "invalid_config_with_duration.txt",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command()
			assert.NotNil(t, cmd)
			cmd.SetArgs(tt.args)
			out := bytes.NewBufferString("")
			cmd.SetOutput(out)

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
