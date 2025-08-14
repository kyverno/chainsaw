package functions

import (
	"errors"
	"fmt"
)

// stable is a utility function that returns the function name as-is.
// Used for stable, non-experimental JMESPath functions.
func stable(in string) string {
	return in
}

// experimental is a utility function that prefixes the function name with "x_".
// Used for experimental JMESPath functions that may change in future versions.
func experimental(in string) string {
	return "x_" + in
}

// getArgAt retrieves an argument at the specified index from the arguments slice.
// Returns an error if the index is out of range.
func getArgAt(arguments []any, index int) (any, error) {
	if index >= len(arguments) {
		return nil, fmt.Errorf("index out of range (%d / %d)", index, len(arguments))
	}
	return arguments[index], nil
}

// getArg is a generic function that retrieves and type-asserts an argument at the specified index.
// Arguments:
// - arguments: The slice of arguments to retrieve from
// - index: The index of the argument to retrieve
// - out: A pointer to store the retrieved value
// Returns an error if the index is out of range or if the type assertion fails.
func getArg[T any](arguments []any, index int, out *T) error {
	arg, err := getArgAt(arguments, index)
	if err != nil {
		return err
	}
	if value, ok := arg.(T); !ok {
		return errors.New("invalid type")
	} else {
		*out = value
		return nil
	}
}
