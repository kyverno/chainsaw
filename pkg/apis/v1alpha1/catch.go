package v1alpha1

// Catch defines actions to be executed on failure.
type Catch struct {
	// Description contains a description of the operation.
	// +optional
	Description string `json:"description,omitempty"`

	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Delete represents a deletion operation.
	// +optional
	Delete *Delete `json:"delete,omitempty"`

	// Describe determines the resource describe collector to execute.
	// +optional
	Describe *Describe `json:"describe,omitempty"`

	// Events determines the events collector to execute.
	// +optional
	Events *Events `json:"events,omitempty"`

	// Get determines the resource get collector to execute.
	// +optional
	Get *Get `json:"get,omitempty"`

	// Lookup determines the resource lookup to execute.
	// +optional
	Lookup *Lookup `json:"lookup,omitempty"`

	// PodLogs determines the pod logs collector to execute.
	// +optional
	PodLogs *PodLogs `json:"podLogs,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`

	// Sleep defines zzzz.
	// +optional
	Sleep *Sleep `json:"sleep,omitempty"`

	// Wait determines the resource wait collector to execute.
	// +optional
	Wait *Wait `json:"wait,omitempty"`
}
