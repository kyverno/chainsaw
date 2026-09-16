package internal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kyverno/chainsaw/pkg/loaders/resource"
	"github.com/kyverno/chainsaw/pkg/logging"
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

// CheckObj attempts to parse the command's stdout as one or more YAML/JSON
// documents (using the same permissive, non-K8s-manifest-only parser used
// for apply/assert resource files), for use as the "obj" argument to
// checks.Check. This lets a check/assert use a plain declarative
// comparison against the command's actual output (e.g. `helm template`
// producing a Kubernetes manifest) instead of only CEL/JMESPath
// expressions against the raw $stdout/$stderr string bindings.
//
// Returns nil if stdout isn't parseable as YAML/JSON at all (e.g. plain
// text output) - checks against $stdout/$stderr as CEL/JMESPath
// expressions keep working exactly as before in that case.
func (c *CommandOutput) CheckObj() any {
	resources, err := resource.Parse(c.Stdout.Bytes(), false)
	if err != nil || len(resources) == 0 {
		return nil
	}
	if len(resources) == 1 {
		return resources[0].UnstructuredContent()
	}
	objs := make([]any, 0, len(resources))
	for i := range resources {
		objs = append(objs, resources[i].UnstructuredContent())
	}
	return objs
}
