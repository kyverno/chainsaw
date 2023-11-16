package v1alpha1

import "github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"

// Script describes a script to run as a part of a test step.
type Script struct {
	// Content defines a shell script (run with "sh -c ...").
	// +optional
	Content string `json:"content,omitempty"`

	// SkipLogOutput removes the output from the command. Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`

	// Check is an assertion tree to validate outcome.
	// +optional
	Check v1alpha1.Any `json:"check,omitempty"`
}
