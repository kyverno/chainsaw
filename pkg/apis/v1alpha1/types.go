package v1alpha1

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Output represents an output binding with a match to determine if the binding must be considered or not.
type Output struct {
	// Binding determines the binding to create when the match succeeds.
	Binding `json:",inline"`

	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`
}

const (
	DefaultApplyTimeout   = 5 * time.Second
	DefaultAssertTimeout  = 30 * time.Second
	DefaultCleanupTimeout = 30 * time.Second
	DefaultDeleteTimeout  = 15 * time.Second
	DefaultErrorTimeout   = 30 * time.Second
	DefaultExecTimeout    = 5 * time.Second
)

// Timeouts contains timeouts per operation.
type Timeouts struct {
	// Apply defines the timeout for the apply operation
	Apply *metav1.Duration `json:"apply,omitempty"`

	// Assert defines the timeout for the assert operation
	Assert *metav1.Duration `json:"assert,omitempty"`

	// Cleanup defines the timeout for the cleanup operation
	Cleanup *metav1.Duration `json:"cleanup,omitempty"`

	// Delete defines the timeout for the delete operation
	Delete *metav1.Duration `json:"delete,omitempty"`

	// Error defines the timeout for the error operation
	Error *metav1.Duration `json:"error,omitempty"`

	// Exec defines the timeout for exec operations
	Exec *metav1.Duration `json:"exec,omitempty"`
}

func durationOrDefault(to *metav1.Duration, def time.Duration) time.Duration {
	if to != nil {
		return to.Duration
	}
	return def
}

func (t Timeouts) ApplyDuration() time.Duration {
	return durationOrDefault(t.Apply, DefaultApplyTimeout)
}

func (t Timeouts) AssertDuration() time.Duration {
	return durationOrDefault(t.Assert, DefaultAssertTimeout)
}

func (t Timeouts) CleanupDuration() time.Duration {
	return durationOrDefault(t.Cleanup, DefaultCleanupTimeout)
}

func (t Timeouts) DeleteDuration() time.Duration {
	return durationOrDefault(t.Delete, DefaultDeleteTimeout)
}

func (t Timeouts) ErrorDuration() time.Duration {
	return durationOrDefault(t.Error, DefaultErrorTimeout)
}

func (t Timeouts) ExecDuration() time.Duration {
	return durationOrDefault(t.Exec, DefaultExecTimeout)
}

func (t Timeouts) Combine(override *Timeouts) Timeouts {
	if override == nil {
		return t
	}
	if override.Apply != nil {
		t.Apply = override.Apply
	}
	if override.Assert != nil {
		t.Assert = override.Assert
	}
	if override.Error != nil {
		t.Error = override.Error
	}
	if override.Delete != nil {
		t.Delete = override.Delete
	}
	if override.Cleanup != nil {
		t.Cleanup = override.Cleanup
	}
	if override.Exec != nil {
		t.Exec = override.Exec
	}
	return t
}
