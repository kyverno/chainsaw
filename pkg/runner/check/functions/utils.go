package functions

import (
	"fmt"
)

func getArg[T any](arguments []any, index int, out *T) error {
	if index >= len(arguments) {
		return fmt.Errorf("index out of range (%d / %d)", index, len(arguments))
	}
	if value, ok := arguments[index].(T); !ok {
		return fmt.Errorf("invalid type")
	} else {
		*out = value
		return nil
	}
}
