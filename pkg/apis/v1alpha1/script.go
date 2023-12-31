package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Script describes a script to run as a part of a test step.
type Script struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Content defines a shell script (run with "sh -c ...").
	// +optional
	Content string `json:"content,omitempty"`

	// SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`

	// Check is an assertion tree to validate the operation outcome.
	// +optional
	Check *Check `json:"check,omitempty"`
}
