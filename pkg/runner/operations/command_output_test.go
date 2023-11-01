package operations

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandOutput(t *testing.T) {
	tests := []struct {
		name        string
		stdout      string
		stderr      string
		expectedOut string
		expectedErr string
	}{
		{
			name:        "normal output and error",
			stdout:      "  hello world  \n",
			stderr:      "  error occurred  \n",
			expectedOut: "hello world",
			expectedErr: "error occurred",
		},
		{
			name:        "empty output and error",
			stdout:      "   \n  ",
			stderr:      "   \n  ",
			expectedOut: "",
			expectedErr: "",
		},
		{
			name:        "no trimming needed",
			stdout:      "hello",
			stderr:      "error",
			expectedOut: "hello",
			expectedErr: "error",
		},
		{
			name:        "newlines only",
			stdout:      "\n\n\n",
			stderr:      "\n\n\n",
			expectedOut: "",
			expectedErr: "",
		},
		{
			name:        "mixed spaces and newlines",
			stdout:      "  \n hello \n world \n  ",
			stderr:      "  \n error \n occurred \n  ",
			expectedOut: "hello \n world",
			expectedErr: "error \n occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := &CommandOutput{
				stdout: *bytes.NewBufferString(tt.stdout),
				stderr: *bytes.NewBufferString(tt.stderr),
			}

			assert.Equal(t, tt.expectedOut, co.Out())
			assert.Equal(t, tt.expectedErr, co.Err())
		})
	}
}
