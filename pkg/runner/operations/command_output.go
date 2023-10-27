package operations

import (
	"bytes"
	"strings"
)

type CommandOutput struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (c *CommandOutput) Out() string {
	return strings.TrimSpace(c.stdout.String())
}

func (c *CommandOutput) Err() string {
	return strings.TrimSpace(c.stderr.String())
}
