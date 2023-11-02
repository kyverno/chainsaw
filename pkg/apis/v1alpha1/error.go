package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Error represents an anticipated error condition that may arise during testing.
// Instead of treating such an error as a test failure, it acknowledges it as expected.
type Error struct {
	// FileRef provides a reference to the file containing the expected error.
	FileRef `json:",inline"`

	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}
