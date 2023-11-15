package internal

import (
	"bytes"
	"strings"
)

type CommandOutput struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func (c *CommandOutput) Out() string {
	return strings.TrimSpace(c.Stdout.String())
}

func (c *CommandOutput) Err() string {
	return strings.TrimSpace(c.Stderr.String())
}
