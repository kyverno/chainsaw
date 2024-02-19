package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Wait specifies how to perform wait operations on resources.
type Wait struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration. Specifies how long to wait for the condition to be met before timing out.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster where the wait operation will be performed (default cluster will be used if not specified).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Resource type on which the wait operation will be applied.
	Resource string `json:"resource"`

	// ResourceName specifies the name of the resource to wait for.
	// +optional
	ResourceName string `json:"resourceName,omitempty"`

	// Condition represents the specific condition to wait for.
	// Example: "Available", "Ready", etc.
	Condition string `json:"condition"`

	// ObjectLabelsSelector determines the selection process of objects based on their labels, applicable when waiting for a condition on multiple resources.
	ObjectLabelsSelector `json:",inline"`

	// IncludeUninitialized indicates whether to include uninitialized resources in the wait operation.
	// +optional
	IncludeUninitialized *bool `json:"includeUninitialized,omitempty"`

	// PollInterval specifies how often to check the condition's status before the timeout is reached.
	// +optional
	PollInterval *metav1.Duration `json:"pollInterval,omitempty"`
}
