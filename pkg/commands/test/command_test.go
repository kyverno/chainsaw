package test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kyverno/chainsaw/pkg/runner"
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
			name:    "default",
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
			out:     filepath.Join(basePath, "with_timeout.txt"),
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
				"--test-dir",
				"dir1,dir2,dir3",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "with_test_dirs.txt"),
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
			out:     filepath.Join(basePath, "with_suppress.txt"),
		},
		{
			name: "skip test with regex",
			args: []string{
				"--include-test-regex",
				"test[4-6]",
				"--exclude-test-regex",
				"test[1-3]",
			},
			wantErr: false,
			out:     filepath.Join(basePath, "with_regex.txt"),
		},
		{
			name: "valid config",
			args: []string{
				"--config",
				filepath.Join(basePath, "config/empty_config.yaml"),
			},
			wantErr: true,
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

func TestPrintSummary(t *testing.T) {
	tests := []struct {
		name        string
		summary     *runner.Summary
		expectedOut string
	}{
		{
			name:        "nil summary",
			summary:     nil,
			expectedOut: "Error: Invalid summary provided.",
		},
		{
			name:        "non-nil summary with values",
			summary:     &runner.Summary{PassedTest: 5, FailedTest: 3},
			expectedOut: "\n--- Test Summary ---\nTotal Tests: 8\nPassed Tests: 5\nFailed Tests: 3\n",
		},
		{
			name:        "non-nil summary with zeros",
			summary:     &runner.Summary{PassedTest: 0, FailedTest: 0},
			expectedOut: "\n--- Test Summary ---\nTotal Tests: 0\nPassed Tests: 0\nFailed Tests: 0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			printSummary(tt.summary, out)

			gotOut := out.String()
			if gotOut != tt.expectedOut {
				t.Errorf("expected output %q, but got %q", tt.expectedOut, gotOut)
			}
		})
	}
}
