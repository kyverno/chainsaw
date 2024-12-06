package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/expressions"
	"github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	identifier = regexp.MustCompile(`^(?:\w+|\(.+\))$`)
	NewAny     = v1alpha1.NewAny
	NewCheck   = v1alpha1.NewAssertionTree
	NewMatch   = v1alpha1.NewAssertionTree
)

type Compiler = v1alpha1.Compiler

// Binding represents a key/value set as a binding in an executing test.
type Binding struct {
	// Name the name of the binding.
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Pattern:=`^(?:\w+|\(.+\))$`
	Name Expression `json:"name"`

	// Compiler defines the default compiler to use when evaluating expressions.
	// +optional
	Compiler *Compiler `json:"compiler,omitempty"`

	// Value value of the binding.
	Value Projection `json:"value"`
}

func (b Binding) CheckName() error {
	if !identifier.MatchString(string(b.Name)) {
		return fmt.Errorf("invalid name %s", b.Name)
	}
	return nil
}

// Check represents a check to be applied on the result of an operation.
type Check = v1alpha1.AssertionTree

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

// Expression defines an expression to be used in string fields.
type Expression string

func (e *Expression) MarshalJSON() ([]byte, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(string(*e))
}

func (e *Expression) UnmarshalJSON(data []byte) error {
	var statement string
	err := json.Unmarshal(data, &statement)
	if err != nil {
		return err
	}
	*e = Expression(statement)
	expression := expressions.Parse(context.TODO(), statement)
	if expression == nil {
		return nil
	}
	if expression.Engine == "" {
		return nil
	}
	parser := parsing.NewParser()
	if _, err := parser.Parse(statement); err != nil {
		return err
	}
	return nil
}

func (e Expression) Value(ctx context.Context, compilers compilers.Compilers, bindings apis.Bindings) (string, error) {
	return expressions.String(ctx, compilers, string(e), bindings)
}

// Format determines the output format (json or yaml).
// +kubebuilder:validation:Type:=string
// +kubebuilder:validation:Pattern:=`^(?:json|yaml|\(.+\))$`
type Format Expression

// Match represents a match condition against an evaluated object.
type Match = v1alpha1.AssertionTree

// ObjectName represents an object namespace and name.
type ObjectName struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace Expression `json:"namespace,omitempty"`

	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name Expression `json:"name,omitempty"`
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
	APIVersion Expression `json:"apiVersion"`

	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind Expression `json:"kind"`
}

// Output represents an output binding with a match to determine if the binding must be considered or not.
type Output struct {
	// Binding determines the binding to create when the match succeeds.
	Binding `json:",inline"`

	// Match defines the matching statement.
	// +optional
	Match *Match `json:"match,omitempty"`
}

// DefaultTimeouts contains defautl timeouts per operation.
type DefaultTimeouts struct {
	// Apply defines the timeout for the apply operation
	// +optional
	// +kubebuilder:default:="5s"
	Apply metav1.Duration `json:"apply"`

	// Assert defines the timeout for the assert operation
	// +optional
	// +kubebuilder:default:="30s"
	Assert metav1.Duration `json:"assert"`

	// Cleanup defines the timeout for the cleanup operation
	// +optional
	// +kubebuilder:default:="30s"
	Cleanup metav1.Duration `json:"cleanup"`

	// Delete defines the timeout for the delete operation
	// +optional
	// +kubebuilder:default:="15s"
	Delete metav1.Duration `json:"delete"`

	// Error defines the timeout for the error operation
	// +optional
	// +kubebuilder:default:="30s"
	Error metav1.Duration `json:"error"`

	// Exec defines the timeout for exec operations
	// +optional
	// +kubebuilder:default:="5s"
	Exec metav1.Duration `json:"exec"`
}

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
