package functions

import (
	"errors"
	"fmt"
)

func stable(in string) string {
	return in
}

func experimental(in string) string {
	return "x_" + in
}

func getArgAt(arguments []any, index int) (any, error) {
	if index >= len(arguments) {
		return nil, fmt.Errorf("index out of range (%d / %d)", index, len(arguments))
	}
	return arguments[index], nil
}

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
