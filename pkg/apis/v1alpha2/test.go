package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

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
// +k8s:conversion-gen=false
type TestSpec struct {
	// Cleanup contains cleanup configuration.
	// +optional
	// +kubebuilder:default:={}
	Cleanup CleanupOptions `json:"cleanup"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters Clusters `json:"clusters,omitempty"`

	// Execution contains tests execution configuration.
	// +optional
	// +kubebuilder:default:={}
	Execution TestExecutionOptions `json:"execution"`

	// Bindings defines additional binding key/values.
	// +optional
	Bindings []Binding `json:"bindings,omitempty"`

	// Deletion contains the global deletion configuration.
	// +optional
	// +kubebuilder:default:={}
	Deletion DeletionOptions `json:"deletion"`

	// Description contains a description of the test.
	// +optional
	Description string `json:"description,omitempty"`

	// Error contains the global error configuration.
	// +optional
	// +kubebuilder:default:={}
	Error ErrorOptions `json:"error"`

	// Namespace contains properties for the namespace to use for tests.
	// +optional
	// +kubebuilder:default:={}
	Namespace NamespaceOptions `json:"namespace"`

	// Steps defining the test.
	Steps []TestStep `json:"steps"`

	// Templating contains the templating config.
	// +optional
	// +kubebuilder:default:={}
	Templating TemplatingOptions `json:"templating"`

	// Timeouts for the test. Overrides the global timeouts set in the Configuration on a per operation basis.
	// +optional
	// +kubebuilder:default:={}
	Timeouts Timeouts `json:"timeouts"`
}

// TestExecutionOptions determines how tests are run.
type TestExecutionOptions struct {
	// Concurrent determines whether the test should run concurrently with other tests.
	// +optional
	// +kubebuilder:default:=true
	Concurrent bool `json:"concurrent"`

	// Skip determines whether the test should skipped.
	// +optional
	Skip bool `json:"skip,omitempty"`

	// TerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.
	// +optional
	TerminationGracePeriod *metav1.Duration `json:"terminationGracePeriod,omitempty"`
}
