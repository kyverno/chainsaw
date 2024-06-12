package v1alpha1

import (
	"fmt"
	"regexp"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

var identifier = regexp.MustCompile(`^(?:\w+|\(.+\))$`)

// Any represents any type.
type Any = v1alpha1.Any

// Binding represents a key/value set as a binding in an executing test.
type Binding struct {
	// Name the name of the binding.
	// +kubebuilder:validation:Pattern=`^(?:\w+|\(.+\))$`
	Name string `json:"name"`

	// Value value of the binding.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	Value Any `json:"value"`
}

func (b Binding) CheckName() error {
	if !identifier.MatchString(b.Name) {
		return fmt.Errorf("invalid name %s", b.Name)
	}
	return nil
}

// Check represents a check to be applied on the result of an operation.
type Check = Any

// Cluster defines cluster config and context.
type Cluster struct {
	// Kubeconfig is the path to the referenced file.
	Kubeconfig string `json:"kubeconfig"`

	// Context is the name of the context to use.
	// +optional
	Context string `json:"context,omitempty"`
}

// Clusters defines a cluster map.
type Clusters map[string]Cluster

// Expectation represents a check to be applied on the result of an operation
// with a match filter to determine if the verification should be considered.
type Expectation struct {
	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`

	// Check defines the verification statement.
	Check Check `json:"check"`
}

// Format determines the output format (json or yaml).
// +kubebuilder:validation:Pattern=`^(?:json|yaml|\(.+\))$`
type Format string

// Match represents a match condition against an evaluated object.
type Match = Any

// ObjectName represents an object namespace and name.
type ObjectName struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty"`
}

// ObjectReference represents one or more objects with a specific apiVersion and kind.
// For a single object name and namespace are used to identify the object.
// For multiple objects use labels.
type ObjectReference struct {
	ObjectType `json:",inline"`
	ObjectName `json:",inline"`

	// Label selector to match objects to delete
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// ObjectType represents a specific apiVersion and kind.
type ObjectType struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion"`

	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`
}
