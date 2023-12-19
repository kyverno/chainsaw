package functions

import (
	"github.com/jmespath-community/go-jmespath/pkg/functions"
	"github.com/kyverno/chainsaw/pkg/client"
)

func GetFunctions(c client.Client) []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: "env",
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpEnv,
	}, {
		Name:      "k8s_client",
		Arguments: []functions.ArgSpec{},
		Handler:   jpKubernetesClient(c),
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
