package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Assert represents a test condition that is expected to hold true
// during the testing process.
type Assert struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ClusterRef is the name of the cluster to run the test on.
	// +optional
	ClusterRef *string `json:"clusterRef,omitempty"`

	// FileRefOrAssert provides a reference to the assertion.
	FileRefOrCheck `json:",inline"`
}
