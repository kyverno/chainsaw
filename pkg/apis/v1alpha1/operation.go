package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Operation struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`

	// Assert represents an assertion to be made. It checks whether the conditions specified in the assertion hold true.
	// +optional
	Assert *Assert `json:"assert,omitempty"`

	// Apply represents resources that should be applied for this test step. This can include things
	// like configuration settings or any other resources that need to be available during the test.
	// +optional
	Apply *Apply `json:"apply,omitempty"`

	// Create represents a creation operation.
	// +optional
	Create *Create `json:"create,omitempty"`

	// Error represents the expected errors for this test step. If any of these errors occur, the test
	// will consider them as expected; otherwise, they will be treated as test failures.
	// +optional
	Error *Error `json:"error,omitempty"`

	// Delete represents a creation operation.
	// +optional
	Delete *Delete `json:"delete,omitempty"`

	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`
}
