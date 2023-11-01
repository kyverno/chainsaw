package discovery

import (
	"regexp"
)

var StepFileName = regexp.MustCompile(`^(\d\d)-(.*)\.(?:yaml|yml)$`)
