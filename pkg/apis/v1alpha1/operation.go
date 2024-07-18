package v1alpha1

// OperationBase defines common elements to all operations.
type OperationBase struct {
	// Description contains a description of the operation.
	// +optional
	Description string `json:"description,omitempty"`

	// ContinueOnError determines whether a test should continue or not in case the operation was not successful.
	// Even if the test continues executing, it will still be reported as failed.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
}

// Operation defines a single operation, only one action is permitted for a given operation.
// +kubebuilder:oneOf:={required:{apply}}
// +kubebuilder:oneOf:={required:{assert}}
// +kubebuilder:oneOf:={required:{command}}
// +kubebuilder:oneOf:={required:{create}}
// +kubebuilder:oneOf:={required:{delete}}
// +kubebuilder:oneOf:={required:{describe}}
// +kubebuilder:oneOf:={required:{error}}
// +kubebuilder:oneOf:={required:{events}}
// +kubebuilder:oneOf:={required:{patch}}
// +kubebuilder:oneOf:={required:{podLogs}}
// +kubebuilder:oneOf:={required:{proxy}}
// +kubebuilder:oneOf:={required:{script}}
// +kubebuilder:oneOf:={required:{sleep}}
// +kubebuilder:oneOf:={required:{update}}
// +kubebuilder:oneOf:={required:{wait}}
type Operation struct {
	// OperationBase defines common elements to all operations.
	// +optional
	OperationBase `json:",inline"`

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

	// Proxy runs a proxy request.
	// +optional
	Proxy *Proxy `json:"proxy,omitempty"`

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

func (o *Operation) Bindings() []Binding {
	switch {
	case o.Apply != nil:
		return o.Apply.Bindings
	case o.Assert != nil:
		return o.Assert.Bindings
	case o.Command != nil:
		return o.Command.Bindings
	case o.Create != nil:
		return o.Create.Bindings
	case o.Delete != nil:
		return o.Delete.Bindings
	case o.Describe != nil:
		return nil
	case o.Error != nil:
		return o.Error.Bindings
	case o.Events != nil:
		return nil
	case o.Get != nil:
		return nil
	case o.Patch != nil:
		return o.Patch.Bindings
	case o.PodLogs != nil:
		return nil
	case o.Proxy != nil:
		return nil
	case o.Script != nil:
		return o.Script.Bindings
	case o.Sleep != nil:
		return nil
	case o.Update != nil:
		return o.Update.Bindings
	case o.Wait != nil:
		return nil
	}
	panic("missing binding operation type handler")
}

func (o *Operation) Outputs() []Output {
	switch {
	case o.Apply != nil:
		return o.Apply.Outputs
	case o.Assert != nil:
		return nil
	case o.Command != nil:
		return o.Command.Outputs
	case o.Create != nil:
		return o.Create.Outputs
	case o.Delete != nil:
		return nil
	case o.Describe != nil:
		return nil
	case o.Error != nil:
		return nil
	case o.Events != nil:
		return nil
	case o.Get != nil:
		return nil
	case o.Patch != nil:
		return o.Patch.Outputs
	case o.PodLogs != nil:
		return nil
	case o.Proxy != nil:
		return o.Proxy.Outputs
	case o.Script != nil:
		return o.Script.Outputs
	case o.Sleep != nil:
		return nil
	case o.Update != nil:
		return o.Update.Outputs
	case o.Wait != nil:
		return nil
	}
	panic("missing output operation type handler")
}
