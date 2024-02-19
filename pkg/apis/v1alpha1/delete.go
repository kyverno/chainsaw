package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Delete is a reference to an object that should be deleted
type Delete struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// ObjectReference determines objects to be deleted.
	ObjectReference `json:"ref"`

	// Expect defines a list of matched checks to validate the operation outcome.
	// +optional
	Expect []Expectation `json:"expect,omitempty"`
}
