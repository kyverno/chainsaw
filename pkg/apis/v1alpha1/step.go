package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// StepTemplate is the resource that contains a step definition.
type StepTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Test step spec.
	Spec StepTemplateSpec `json:"spec"`
}

// StepTemplateSpec defines the spec of a step template.
type StepTemplateSpec struct {
	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Try defines what the step will try to execute.
	// +kubebuilder:validation:MinItems:=1
	Try []Operation `json:"try"`

	// Catch defines what the step will execute when an error happens.
	// +optional
	Catch []CatchFinally `json:"catch,omitempty"`

	// Finally defines what the step will execute after the step is terminated.
	// +optional
	Finally []CatchFinally `json:"finally,omitempty"`

	// Cleanup defines what will be executed after the test is terminated.
	// +optional
	Cleanup []CatchFinally `json:"cleanup,omitempty"`
}

// TestStep contains the test step definition used in a test spec.
// +kubebuilder:not:={required:{use},anyOf:{{required:{try}},{required:{catch}},{required:{finally}},{required:{cleanup}}}}
type TestStep struct {
	// Name of the step.
	// +optional
	Name string `json:"name,omitempty"`

	// Use defines a reference to a step template.
	Use *Use `json:"use,omitempty"`

	// TestStepSpec of the step.
	TestStepSpec `json:",inline"`
}

// Use defines a reference to a step template.
type Use struct {
	// Template references a step template.
	Template string `json:"template"`

	// With defines arguments passed to the step template.
	// +optional
	// +kubebuilder:default:={}
	With With `json:"with"`
}

// With defines arguments passed to step templates.
type With struct {
	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`
}

// TestStepSpec defines the desired state and behavior for each test step.
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

	// Cluster defines the target cluster (will be inherited if not specified).
	// +optional
	Cluster *string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// Skip determines whether the step should be skipped. Can be a boolean or a template expression.
	// +optional
	Skip *string `json:"skip,omitempty"`

	// SkipDelete determines whether the resources created by the step should be deleted after the test step is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// Compiler defines the default compiler to use when evaluating expressions.
	// +optional
	Compiler *Compiler `json:"compiler,omitempty"`

	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Try defines what the step will try to execute.
	// +kubebuilder:validation:MinItems:=1
	// +optional
	Try []Operation `json:"try"`

	// Catch defines what the step will execute when an error happens.
	// +optional
	Catch []CatchFinally `json:"catch,omitempty"`

	// Finally defines what the step will execute after the step is terminated.
	// +optional
	Finally []CatchFinally `json:"finally,omitempty"`

	// Cleanup defines what will be executed after the test is terminated.
	// +optional
	Cleanup []CatchFinally `json:"cleanup,omitempty"`
}
