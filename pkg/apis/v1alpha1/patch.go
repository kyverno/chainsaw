package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Patch represents a set of resources that should be patched.
// If a resource doesn't exist yet in the cluster it will fail.
type Patch struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Outputs defines output bindings.
	// +optional
	Outputs []Output `json:"outputs,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters map[string]Cluster `json:"clusters,omitempty"`

	// FileRefOrResource provides a reference to the file containing the resources to be patched.
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
