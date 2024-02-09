package v1alpha1

import (
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Patch struct {
	// Match defines the matching statement.
	// +optional
	Match *Check `json:"match,omitempty"`

	// Check defines the verification statement.
	Patch v1alpha1.Any `json:"patch"`
}

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// FileRefOrResource provides a reference to the resources to be applied.
	FileRefOrResource `json:",inline"`

	Patches []Patch `json:"patches,omitempty"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Expect defines a list of matched checks to validate the operation outcome.
	// +optional
	Expect []Expectation `json:"expect,omitempty"`
}
