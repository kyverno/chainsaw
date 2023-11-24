package v1alpha1

import (
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
)

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// FileRefOrResource provides a reference to the file containing the resources to be applied.
	FileRefOrResource `json:",inline"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Check is an assertion tree to validate outcome.
	// +optional
	Check *v1alpha1.Any `json:"check,omitempty"`
}
