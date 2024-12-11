package internal

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestCommandOutput(t *testing.T) {
	tests := []struct {
		name        string
		stdout      string
		stderr      string
		expectedOut string
		expectedErr string
	}{{
		name:        "normal output and error",
		stdout:      "  hello world  \n",
		stderr:      "  error occurred  \n",
		expectedOut: "  hello world  \n",
		expectedErr: "  error occurred  \n",
	}, {
		name:        "empty output and error",
		stdout:      "   \n  ",
		stderr:      "   \n  ",
		expectedOut: "   \n  ",
		expectedErr: "   \n  ",
	}, {
		name:        "no trimming needed",
		stdout:      "hello",
		stderr:      "error",
		expectedOut: "hello",
		expectedErr: "error",
	}, {
		name:        "newlines only",
		stdout:      "\n\n\n",
		stderr:      "\n\n\n",
		expectedOut: "\n\n\n",
		expectedErr: "\n\n\n",
	}, {
		name:        "mixed spaces and newlines",
		stdout:      "  \n hello \n world \n  ",
		stderr:      "  \n error \n occurred \n  ",
		expectedOut: "  \n hello \n world \n  ",
		expectedErr: "  \n error \n occurred \n  ",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := &CommandOutput{
				Stdout: *bytes.NewBufferString(tt.stdout),
				Stderr: *bytes.NewBufferString(tt.stderr),
			}
			assert.Equal(t, tt.expectedOut, co.Out())
			assert.Equal(t, tt.expectedErr, co.Err())
		})
	}
}

func TestCommandOutput_Sections(t *testing.T) {
	var o, e bytes.Buffer
	o.WriteString("out")
	e.WriteString("err")
	tests := []struct {
		name   string
		Stdout bytes.Buffer
		Stderr bytes.Buffer
		want   []fmt.Stringer
	}{{
		name: "none",
		want: nil,
	}, {
		name:   "out",
		Stdout: o,
		want: []fmt.Stringer{
			logging.Section("STDOUT", o.String()),
		},
	}, {
		name:   "err",
		Stderr: e,
		want: []fmt.Stringer{
			logging.Section("STDERR", e.String()),
		},
	}, {
		name:   "both",
		Stdout: o,
		Stderr: e,
		want: []fmt.Stringer{
			logging.Section("STDOUT", o.String()),
			logging.Section("STDERR", e.String()),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandOutput{
				Stdout: tt.Stdout,
				Stderr: tt.Stderr,
			}
			got := c.Sections()
			assert.Equal(t, tt.want, got)
		})
	}
}
