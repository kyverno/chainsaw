package v1alpha1

// Operation defines a single operation, only one action is permitted for a given operation.
type Operation struct {
	// Description contains a description of the operation.
	// +optional
	Description string `json:"description,omitempty"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`

	// Apply represents resources that should be applied for this test step. This can include things
	// like configuration settings or any other resources that need to be available during the test.
	// +optional
	Apply *Apply `json:"apply,omitempty"`

	// Assert represents an assertion to be made. It checks whether the conditions specified in the assertion hold true.
	// +optional
	Assert *Assert `json:"assert,omitempty"`

	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Create represents a creation operation.
	// +optional
	Create *Create `json:"create,omitempty"`

	// Delete represents a deletion operation.
	// +optional
	Delete *Delete `json:"delete,omitempty"`

	// Error represents the expected errors for this test step. If any of these errors occur, the test
	// will consider them as expected; otherwise, they will be treated as test failures.
	// +optional
	Error *Error `json:"error,omitempty"`

	// Lookup determines the resource lookup to execute.
	// +optional
	Lookup *Lookup `json:"lookup,omitempty"`

	// Patch represents a patch operation.
	// +optional
	Patch *Patch `json:"patch,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`

	// Sleep defines zzzz.
	// +optional
	Sleep *Sleep `json:"sleep,omitempty"`

	// Update represents an update operation.
	// +optional
	Update *Update `json:"update,omitempty"`

	// Wait determines the resource wait collector to execute.
	// +optional
	Wait *Wait `json:"wait,omitempty"`
}
