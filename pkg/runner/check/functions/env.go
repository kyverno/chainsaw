package functions

import (
	"errors"
	"fmt"
	"os"
)

func jpEnv(arguments []any) (any, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("invalid number of arguments, expected %d, got %d", 1, len(arguments))
	}
	if key, ok := arguments[0].(string); !ok {
		return nil, errors.New("invalid type, first argument must be a string")
	} else {
		return os.Getenv(key), nil
	}
}
