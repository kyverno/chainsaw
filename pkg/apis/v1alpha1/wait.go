package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Wait specifies how to perform wait operations on resources.
type Wait struct {
	// Timeout for the operation. Specifies how long to wait for the condition to be met before timing out.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster where the wait operation will be performed (default cluster will be used if not specified).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// ResourceReference referenced resource type.
	ResourceReference `json:",inline"`

	// ObjectLabelsSelector determines the selection process of referenced objects.
	ObjectLabelsSelector `json:",inline"`

	// For specifies the condition to wait for.
	For `json:"for"`

	// Format determines the output format (json or yaml).
	// +optional
	Format Format `json:"format,omitempty"`
}
