package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// WaitType represents the type of wait operation.
type WaitType string

// Define constants for WaitType values
const (
	WaitTypeDelete    WaitType = "delete"
	WaitTypeCondition WaitType = "condition"
)

// Wait specifies how to perform wait operations on resources.
type Wait struct {
	// Timeout for the operation. Specifies how long to wait for the condition to be met before timing out.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster where the wait operation will be performed (default cluster will be used if not specified).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Resource type on which the wait operation will be applied.
	Resource string `json:"resource"`

	// ObjectLabelsSelector determines the selection process of referenced objects.
	ObjectLabelsSelector `json:",inline"`

	// Condition represents the specific condition to wait for.
	// Example: "Available", "Ready", etc.
	// +optional
	Condition string `json:"condition,omitempty"`

	// WaitType specifies the type of wait operation: "condition" or "delete".
	// +kubebuilder:validation:Enum=delete;condition
	WaitType WaitType `json:"waitType"`

	// AllNamespaces indicates whether to wait for resources in all namespaces.
	// +optional
	AllNamespaces bool `json:"allNamespaces,omitempty"`
}
