package functions

import (
	"os"
)

func jpEnv(arguments []any) (any, error) {
	var key string
	if err := getArg(arguments, 0, &key); err != nil {
		return nil, err
	}
	return os.Getenv(key), nil
}
