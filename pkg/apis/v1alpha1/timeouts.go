package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Timeouts contains timeouts per operation.
type Timeouts struct {
	// Apply defines the timeout for the apply operation
	Apply *metav1.Duration `json:"apply,omitempty"`

	// Assert defines the timeout for the assert operation
	Assert *metav1.Duration `json:"assert,omitempty"`

	// Error defines the timeout for the error operation
	Error *metav1.Duration `json:"error,omitempty"`

	// Delete defines the timeout for the delete operation
	Delete *metav1.Duration `json:"delete,omitempty"`

	// Cleanup defines the timeout for the cleanup operation
	Cleanup *metav1.Duration `json:"cleanup,omitempty"`

	// Exec defines the timeout for exec operations
	Exec *metav1.Duration `json:"exec,omitempty"`
}
