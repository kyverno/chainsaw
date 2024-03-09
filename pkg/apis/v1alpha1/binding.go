package v1alpha1

import (
	"fmt"
	"regexp"
)

var identifier = regexp.MustCompile(`^(?:\w+|\(.+\))$`)

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
