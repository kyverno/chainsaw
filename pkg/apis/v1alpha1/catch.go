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

	// Describe determines the resource describe collector to execute.
	// +optional
	Describe *Describe `json:"describe,omitempty"`

	// Wait determines the resource wait collector to execute.
	// +optional
	Wait *Wait `json:"wait,omitempty"`

	// Get determines the resource get collector to execute.
	// +optional
	Get *Get `json:"get,omitempty"`

	// Delete represents a deletion operation.
	// +optional
	Delete *Delete `json:"delete,omitempty"`

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

func (c *Catch) Bindings() []Binding {
	switch {
	case c.Command != nil:
		return c.Command.Bindings
	case c.Delete != nil:
		return c.Delete.Bindings
	case c.Describe != nil:
		return nil
	case c.Events != nil:
		return nil
	case c.Get != nil:
		return nil
	case c.PodLogs != nil:
		return nil
	case c.Script != nil:
		return c.Script.Bindings
	case c.Sleep != nil:
		return nil
	case c.Wait != nil:
		return nil
	}
	panic("missing binding operation type handler")
}

func (c *Catch) Outputs() []Output {
	switch {
	case c.Command != nil:
		return c.Command.Outputs
	case c.Delete != nil:
		return nil
	case c.Describe != nil:
		return nil
	case c.Events != nil:
		return nil
	case c.Get != nil:
		return nil
	case c.PodLogs != nil:
		return nil
	case c.Script != nil:
		return c.Script.Outputs
	case c.Sleep != nil:
		return nil
	case c.Wait != nil:
		return nil
	}
	panic("missing output operation type handler")
}
