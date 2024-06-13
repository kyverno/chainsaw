package v1alpha2

// Operation defines operation elements.
// +k8s:conversion-gen=false
type Operation struct {
	OperationAction   `json:",inline"`
	OperationBindings `json:",inline"`
	OperationClusters `json:",inline"`
	OperationOutputs  `json:",inline"`

	// Description contains a description of the operation.
	// +optional
	Description string `json:"description,omitempty"`
}

// OperationAction defines an operation action, only one action should be specified per operation.
type OperationAction struct {
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

	// Describe determines the resource describe collector to execute.
	// +optional
	Describe *Describe `json:"describe,omitempty"`

	// Error represents the expected errors for this test step. If any of these errors occur, the test
	// will consider them as expected; otherwise, they will be treated as test failures.
	// +optional
	Error *Error `json:"error,omitempty"`

	// Events determines the events collector to execute.
	// +optional
	Events *Events `json:"events,omitempty"`

	// Get determines the resource get collector to execute.
	// +optional
	Get *Get `json:"get,omitempty"`

	// Patch represents a patch operation.
	// +optional
	Patch *Patch `json:"patch,omitempty"`

	// PodLogs determines the pod logs collector to execute.
	// +optional
	PodLogs *PodLogs `json:"podLogs,omitempty"`

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

// OperationBindings contains bindings options for an operation.
type OperationBindings struct {
	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`
}

// OperationClusters contains clusters options for an operation.
type OperationClusters struct {
	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`
}

// OperationOutputs contains outputs options for an operation.
type OperationOutputs struct {
	// Outputs defines output bindings.
	// +optional
	Outputs []Output `json:"outputs,omitempty"`
}

// TryOperation defines operation elements.
// +k8s:conversion-gen=false
type TryOperation struct {
	Operation `json:",inline"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}
