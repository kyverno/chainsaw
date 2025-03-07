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
	path := "../../../.temp"
	basePath := "../../../testdata/commands/test"
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
		err     string
	}{{
		name: "help",
		args: []string{
			"--help",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "help.txt"),
	}, {
		name:    "default",
		args:    []string{},
		wantErr: false,
		out:     filepath.Join(basePath, "default.txt"),
	}, {
		name: "with apply timeout",
		args: []string{
			"--apply-timeout",
			"10s",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "with_timeout.txt"),
	}, {
		name: "with repeat count",
		args: []string{
			"--repeat-count",
			"3",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "with_repeat_count.txt"),
	}, {
		name: "with no-warnings flag",
		args: []string{
			"--no-warnings",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "with_no_warnings.txt"),
	}, {
		name: "invalid timeout",
		args: []string{
			"--timeout",
			"invalid",
		},
		wantErr: true,
	}, {
		name: "test dirs specified",
		args: []string{
			"--test-dir",
			"..",
			"--test-dir",
			".",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "with_test_dirs.txt"),
	}, {
		name: "nonexistent config file",
		args: []string{
			"--config",
			"nonexistent.yaml",
		},
		wantErr: true,
	}, {
		name: "skip test with regex",
		args: []string{
			"--include-test-regex",
			"test[4-6]",
			"--exclude-test-regex",
			"test[1-3]",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "with_regex.txt"),
	}, {
		name: "empty config",
		args: []string{
			"--config",
			filepath.Join(basePath, "config/empty_config.yaml"),
		},
		wantErr: true,
	}, {
		name: "nonexistent config",
		args: []string{
			"--config",
			filepath.Join(basePath, "config/nonexistent_config.yaml"),
		},
		wantErr: true,
	}, {
		name: "misformatted config",
		args: []string{
			"--config",
			filepath.Join(basePath, "config/wrong_format_config.yaml"),
		},
		wantErr: true,
		out:     filepath.Join(basePath, "wrong_format_config.txt"),
	}, {
		name: "wrong kind in config",
		args: []string{
			"--config",
			filepath.Join(basePath, "config/wrong_kind_config.yaml"),
		},
		wantErr: true,
		out:     filepath.Join(basePath, "wrong_kind_config.txt"),
		err:     filepath.Join(basePath, "wrong_kind_config_err.txt"),
	}, {
		name: "config with all fields",
		args: []string{
			"--config",
			filepath.Join(basePath, "config/config_all_fields.yaml"),
			"--report-path",
			path,
		},
		wantErr: false,
		out:     filepath.Join(basePath, "config_all_fields.txt"),
	}, {
		name: "all flags",
		args: []string{
			"--test-file=custom-test.yaml",
			"--apply-timeout=100s",
			"--assert-timeout=100s",
			"--error-timeout=100s",
			"--delete-timeout=100s",
			"--cleanup-timeout=100s",
			"--exec-timeout=100s",
			"--test-dir=.",
			"--skip-delete=false",
			"--fail-fast=false",
			"--parallel=24",
			"--repeat-count=12",
			"--report-format=XML",
			"--report-path",
			path,
			"--report-name=foo",
			"--namespace=bar",
			"--full-name=true",
			"--include-test-regex=^.*$",
			"--exclude-test-regex=^.*$",
			"--force-termination-grace-period=5s",
		},
		wantErr: false,
		out:     filepath.Join(basePath, "all_flags.txt"),
	}}
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
