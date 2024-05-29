package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Describe defines how to describe resources.
type Describe struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// ResourceReference referenced resource type.
	ResourceReference `json:",inline"`

	// ObjectLabelsSelector determines the selection process of referenced objects.
	// +optional
	ObjectLabelsSelector `json:",inline"`

	// Show Events indicates whether to include related events.
	// +optional
	ShowEvents *bool `json:"showEvents,omitempty"`
}
