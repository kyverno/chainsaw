package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Command describes a command to run as a part of a test step.
type Command struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Entrypoint is the command entry point to run.
	Entrypoint string `json:"entrypoint"`

	// Args is the command arguments.
	// +optional
	Args []string `json:"args,omitempty"`

	// SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`

	// Check is an assertion tree to validate the operation outcome.
	// +optional
	Check *Check `json:"check,omitempty"`
}
