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
	}, {
		Name: "k8s_list",
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpKubernetesList,
	}}
}
