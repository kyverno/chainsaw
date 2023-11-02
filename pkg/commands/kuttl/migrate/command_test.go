package migrate

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/kyverno/chainsaw/pkg/commands/root"
	"github.com/stretchr/testify/assert"
)

func Test_Execute(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		out     string
	}{{
		name: "help",
		args: []string{
			"migrate",
			"--help",
		},
		out:     "../../../../testdata/commands/kuttl/migrate/help.txt",
		wantErr: false,
	}, {
		name: "migrate",
		args: []string{
			"migrate",
			"../../../../testdata/kuttl",
		},
		out:     "../../../../testdata/commands/kuttl/migrate/out.txt",
		wantErr: false,
	}, {
		name: "migrate save",
		args: []string{
			"migrate",
			"../../../../testdata/kuttl",
			"--save",
		},
		out:     "../../../../testdata/commands/kuttl/migrate/out-save.txt",
		wantErr: false,
	}, {
		name: "unknow flag",
		args: []string{
			"migrate",
			"--foo",
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
