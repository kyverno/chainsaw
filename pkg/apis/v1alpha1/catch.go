package v1alpha1

// Catch defines actions to be executed on failure.
type Catch struct {
	// Description contains a description of the operation.
	// +optional
	Description string `json:"description,omitempty"`

	// PodLogs determines the pod logs collector to execute.
	// +optional
	PodLogs *PodLogs `json:"podLogs,omitempty"`

	// Events determines the events collector to execute.
	// +optional
	Events *Events `json:"events,omitempty"`

	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`

	// Sleep defines zzzz.
	// +optional
	Sleep *Sleep `json:"sleep,omitempty"`
}
