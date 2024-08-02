package internal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kyverno/chainsaw/pkg/engine/logging"
)

type CommandOutput struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func (c *CommandOutput) Out() string {
	return c.Stdout.String()
}

func (c *CommandOutput) Err() string {
	return c.Stderr.String()
}

func (c *CommandOutput) Sections() []fmt.Stringer {
	var sections []fmt.Stringer
	o := strings.TrimSpace(c.Out())
	e := strings.TrimSpace(c.Err())
	if o != "" {
		sections = append(sections, logging.Section("STDOUT", o))
	}
	if e != "" {
		sections = append(sections, logging.Section("STDERR", e))
	}
	return sections
}
