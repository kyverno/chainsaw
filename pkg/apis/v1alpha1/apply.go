package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Apply represents a set of configurations or resources that
// should be applied during testing.
type Apply struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// FileRefOrResource provides a reference to the resources to be applied.
	FileRefOrResource `json:",inline"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// DryRun determines whether the file should be applied in dry run mode.
	// +optional
	DryRun *bool `json:"dryRun,omitempty"`

	// Expect defines a list of matched checks to validate the operation outcome.
	// +optional
	Expect []Expectation `json:"expect,omitempty"`
}
