package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Exec describes a command or script.
type Exec struct {
	// Timeout for the command. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`
}
