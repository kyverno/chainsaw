package v1alpha1

import (
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
)

// Command describes a command to run as a part of a test step.
type Command struct {
	// Entrypoint is the command entry point to run.
	Entrypoint string `json:"entrypoint"`

	// Args is the command arguments.
	// +optional
	Args []string `json:"args,omitempty"`

	// SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`

	// Check is an assertion tree to validate outcome.
	// +optional
	Check *v1alpha1.Any `json:"check,omitempty"`
}
