package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// Test is the resource that contains a test definition.
type Test struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Test spec.
	Spec TestSpec `json:"spec"`
}

// TestSpec contains the test spec.
type TestSpec struct {
	// Description contains a description of the test.
	// +optional
	Description string `json:"description,omitempty"`

	// FailFast determines whether the test should stop upon encountering the first failure.
	// +optional
	FailFast *bool `json:"failFast,omitempty"`

	// Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`

	// Cluster defines the target cluster (will be inherited if not specified).
	// +optional
	Cluster *string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// Skip determines whether the test should skipped.
	// +optional
	Skip *bool `json:"skip,omitempty"`

	// Concurrent determines whether the test should run concurrently with other tests.
	// +optional
	Concurrent *bool `json:"concurrent,omitempty"`

	// SkipDelete determines whether the resources created by the test should be deleted after the test is executed.
	// +optional
	SkipDelete *bool `json:"skipDelete,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// Namespace determines whether the test should run in a random ephemeral namespace or not.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NamespaceTemplate defines a template to create the test namespace.
	// +optional
	NamespaceTemplate *Projection `json:"namespaceTemplate,omitempty"`

	// Scenarios defines test scenarios.
	// +optional
	Scenarios []Scenario `json:"scenarios,omitempty"`

	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Steps defining the test.
	Steps []TestStep `json:"steps"`

	// Catch defines what the steps will execute when an error happens.
	// This will be combined with catch handlers defined at the step level.
	// +optional
	Catch []CatchFinally `json:"catch,omitempty"`

	// ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.
	// +optional
	ForceTerminationGracePeriod *metav1.Duration `json:"forceTerminationGracePeriod,omitempty"`

	// DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.
	// +optional
	DelayBeforeCleanup *metav1.Duration `json:"delayBeforeCleanup,omitempty"`

	// DeletionPropagationPolicy decides if a deletion will propagate to the dependents of
	// the object, and how the garbage collector will handle the propagation.
	// Overrides the deletion propagation policy set in the Configuration.
	// +optional
	// +kubebuilder:validation:Enum:=Orphan;Background;Foreground
	DeletionPropagationPolicy *metav1.DeletionPropagation `json:"deletionPropagationPolicy,omitempty"`
}

// Scenario defines per scenario bindings.
type Scenario struct {
	// Bindings defines binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`
}
