package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// For specifies the condition to wait for.
type For struct {
	// Deletion specifies parameters for waiting on a resource's deletion.
	// +optional
	Deletion *Deletion `json:"delete,omitempty"`
	// Condition specifies the condition to wait for.
	// +optional
	Condition *Condition `json:"condition,omitempty"`
}

// Deletion represents parameters for waiting on a resource's deletion.
type Deletion struct{}

// Condition represents parameters for waiting on a specific condition of a resource.
type Condition struct {
	// Name is the specific condition to wait for, e.g., "Available", "Ready".
	Name string `json:"name"`
	// Value defines the specific condition status to wait for, e.g., "True", "False".
	// +optional
	Value *string `json:"value,omitempty"`
}

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

	// For specifies the condition to wait for.
	For `json:"for"`

	// OutputFormat specifies the output format for the wait operation.
	// +optional
	OutputFormat *Format `json:"output,omitempty"`
}
