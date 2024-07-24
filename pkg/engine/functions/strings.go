package functions

import (
	"strings"
)

func jpTrimSpace(arguments []any) (any, error) {
	var in string
	if err := getArg(arguments, 0, &in); err != nil {
		return nil, err
	}
	return strings.TrimSpace(in), nil
}
