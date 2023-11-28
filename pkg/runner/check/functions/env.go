package functions

import (
	"errors"
	"os"
)

func jpEnv(arguments []interface{}) (interface{}, error) {
	if key, ok := arguments[0].(string); !ok {
		return nil, errors.New("invalid type, first argument must be a string")
	} else {
		return os.Getenv(key), nil
	}
}
