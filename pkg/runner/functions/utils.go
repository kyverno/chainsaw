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

func getArg[T any](arguments []any, index int, out *T) error {
	if index >= len(arguments) {
		return fmt.Errorf("index out of range (%d / %d)", index, len(arguments))
	}
	if value, ok := arguments[index].(T); !ok {
		return errors.New("invalid type")
	} else {
		*out = value
		return nil
	}
}
