package v1alpha1

type Operation struct {
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

	// Exec provides a command or script that should be executed as a part of this test step.
	// +optional
	Exec *ExecOperation `json:"exec,omitempty"`
}
