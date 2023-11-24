package v1alpha1

import (
	"github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
)

// Create represents a set of resources that should be created.
// If a resource already exists in the cluster it will fail.
type Create struct {
	// FileRefOrResource provides a reference to the file containing the resources to be created.
	FileRefOrResource `json:",inline"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Check is an assertion tree to validate outcome.
	// +optional
	Check *v1alpha1.Any `json:"check,omitempty"`
}
