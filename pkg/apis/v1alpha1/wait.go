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

	// Name specifies the name of the resource to wait for.
	// If empty, wait for all resources of the type.
	// +optional
	Name string `json:"name,omitempty"`

	// Namespace specifies the namespace of the resources to wait for.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Selector to filter resources based on label.
	// +optional
	Selector string `json:"selector,omitempty"`

	// Condition represents the specific condition to wait for.
	// Example: "Available", "Ready", etc.
	Condition string `json:"condition"`

	// AllNamespaces indicates whether to wait for resources in all namespaces.
	// +optional
	AllNamespaces bool `json:"allNamespaces,omitempty"`
}
