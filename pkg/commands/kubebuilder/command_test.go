package kubebuilder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubebuilderCommand(t *testing.T) {
	cmd := Command()
	assert.Equal(t, "kubebuilder", cmd.Use)
	assert.NotNil(t, cmd)
}

func TestScaffoldCommand_stdout(t *testing.T) {
	cmd := Command()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"scaffold", "my-operator-test", "--kind", "Foo", "--api-version", "foo.io/v1alpha1"})
	err := cmd.Execute()
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "chainsaw.kyverno.io")
	assert.Contains(t, out, "my-operator-test")
}
