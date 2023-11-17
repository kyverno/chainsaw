package internal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kyverno/chainsaw/pkg/runner/logging"
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

func (c *CommandOutput) Sections() []fmt.Stringer {
	var sections []fmt.Stringer
	o := c.Out()
	e := c.Err()
	if o != "" {
		sections = append(sections, logging.Section("STDOUT", o))
	}
	if e != "" {
		sections = append(sections, logging.Section("STDERR", e))
	}
	return sections
}
