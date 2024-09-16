package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestStep contains the test step definition used in a test spec.
// +k8s:conversion-gen=false
type TestStep struct {
	// Name of the step.
	// +optional
	Name string `json:"name,omitempty"`

	// TestStepSpec of the step.
	TestStepSpec `json:",inline"`
}

// TestStepSpec defines the desired state and behavior for each test step.
// +k8s:conversion-gen=false
type TestStepSpec struct {
	// Description contains a description of the test step.
	// +optional
	Description string `json:"description,omitempty"`

	// Timeouts for the test step. Overrides the global timeouts set in the Configuration and the timeouts eventually set in the Test.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`

	// DeletionPropagationPolicy decides if a deletion will propagate to the dependents of
	// the object, and how the garbage collector will handle the propagation.
	// Overrides the deletion propagation policy set in both the Configuration and the Test.
	// +optional
	// +kubebuilder:validation:Enum:=Orphan;Background;Foreground
	DeletionPropagationPolicy *metav1.DeletionPropagation `json:"deletionPropagationPolicy,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Try defines what the step will try to execute.
	// +kubebuilder:validation:MinItems:=1
	Try []TryOperation `json:"try"`

	// Catch defines what the step will execute when an error happens.
	// +optional
	Catch []Operation `json:"catch,omitempty"`

	// Finally defines what the step will execute after the step is terminated.
	// +optional
	Finally []Operation `json:"finally,omitempty"`

	// Cleanup defines what will be executed after the test is terminated.
	// +optional
	Cleanup []Operation `json:"cleanup,omitempty"`
}
