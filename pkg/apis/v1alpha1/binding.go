package v1alpha1

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	identifier     = regexp.MustCompile(`^\w+$`)
	forbiddenNames = []string{"namespace", "client", "config", "error", "values", "stdout", "stderr"}
	forbidden      = sets.New(forbiddenNames...)
)

// Binding represents a key/value set as a binding in an executing test.
type Binding struct {
	// Name the name of the binding.
	Name string `json:"name"`

	// Value value of the binding.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	Value Any `json:"value"`
}

func (b Binding) CheckName() error {
	return CheckBindingName(b.Name)
}

func (b Binding) CheckEnvName() error {
	return CheckBindingEnvName(b.Name)
}

func CheckBindingName(name string) error {
	if forbidden.Has(name) {
		return fmt.Errorf("name is forbidden (%s), it must not be (%s)", name, strings.Join(forbiddenNames, ", "))
	}
	if !identifier.MatchString(name) {
		return fmt.Errorf("invalid name %s", name)
	}
	return nil
}

func CheckBindingEnvName(name string) error {
	if !identifier.MatchString(name) {
		return fmt.Errorf("invalid name %s", name)
	}
	return nil
}
