package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Create represents a set of resources that should be created.
// If a resource already exists in the cluster it will fail.
type Create struct {
	// FileRef provides a reference to the file containing the
	FileRef `json:",inline"`

	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ShouldFail determines whether applying the file is expected to fail.
	// +optional
	ShouldFail *bool `json:"shouldFail,omitempty"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}
