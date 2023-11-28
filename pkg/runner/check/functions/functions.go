package functions

import (
	"github.com/jmespath-community/go-jmespath/pkg/functions"
)

func GetFunctions() []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: "env",
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpEnv,
	}}
}
