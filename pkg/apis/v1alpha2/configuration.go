package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// Configuration is the resource that contains the configuration used to run tests.
type Configuration struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Configuration spec.
	Spec ConfigurationSpec `json:"spec"`
}

// ConfigurationSpec contains the configuration used to run tests.
// +k8s:conversion-gen=false
type ConfigurationSpec struct {
	// Catch defines what the tests steps will execute when an error happens.
	// This will be combined with catch handlers defined at the test and step levels.
	// +optional
	Catch []Catch `json:"catch,omitempty"`

	// Cleanup contains cleanup configuration.
	// +optional
	Cleanup *Cleanup `json:"cleanup,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters map[string]Cluster `json:"clusters,omitempty"`

	// Discovery contains tests discovery configuration.
	// +optional
	// +kubebuilder:default:={}
	Discovery Discovery `json:"discovery"`

	// Execution contains tests execution configuration.
	// +optional
	Execution *Execution `json:"execution,omitempty"`

	// Namespace contains properties for the namespace to use for tests.
	// +optional
	Namespace *Namespace `json:"namespace,omitempty"`

	// DeletionPropagationPolicy decides if a deletion will propagate to the dependents of
	// the object, and how the garbage collector will handle the propagation.
	// +optional
	// +kubebuilder:validation:Enum:=Orphan;Background;Foreground
	DeletionPropagationPolicy *metav1.DeletionPropagation `json:"deletionPropagationPolicy,omitempty"`

	// Report contains properties for the report.
	// +optional
	Report *Report `json:"report,omitempty"`

	// Templating contains the templating config.
	// +optional
	Templating *Templating `json:"templating,omitempty"`

	// Global timeouts configuration. Applies to all tests/test steps if not overridden.
	// +optional
	Timeouts Timeouts `json:"timeouts"`
}
