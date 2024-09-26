package functions

import (
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/functions"
)

var (
	// stable functions
	env       = stable("env")
	trimSpace = stable("trim_space")
	asString  = stable("as_string")
	// experimental functions
	k8sGet            = experimental("k8s_get")
	k8sList           = experimental("k8s_list")
	k8sExists         = experimental("k8s_exists")
	k8sResourceExists = experimental("k8s_resource_exists")
	k8sServerVersion  = experimental("k8s_server_version")
	metricsDecode     = experimental("metrics_decode")
)

func GetFunctions() []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: env,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler:     jpEnv,
		Description: "Returns the value of the environment variable passed in argument.",
	}, {
		Name: k8sGet,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler:     jpKubernetesGet,
		Description: "Gets a resource from a Kubernetes cluster.",
	}, {
		Name: k8sList,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler:     jpKubernetesList,
		Description: "Lists resources from a Kubernetes cluster.",
	}, {
		Name: k8sExists,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler:     jpKubernetesExists,
		Description: "Checks if a given resource exists in a Kubernetes cluster.",
	}, {
		Name: k8sResourceExists,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler:     jpKubernetesResourceExists,
		Description: "Checks if a given resource type is available in a Kubernetes cluster.",
	}, {
		Name: k8sServerVersion,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
		},
		Handler:     jpKubernetesServerVersion,
		Description: "Returns the version of a Kubernetes cluster.",
	}, {
		Name: metricsDecode,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler:     jpMetricsDecode,
		Description: "Decodes metrics in the Prometheus text format.",
	}, {
		Name: trimSpace,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler:     jpTrimSpace,
		Description: "Trims leading and trailing spaces from the string passed in argument.",
	}, {
		Name: asString,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
		},
		Handler: func(arguments []any) (any, error) {
			in, err := getArgAt(arguments, 0)
			if err != nil {
				return nil, err
			}
			if in != nil {
				if in, ok := in.(string); ok {
					return in, nil
				}
				if reflect.ValueOf(in).Kind() == reflect.String {
					return fmt.Sprint(in), nil
				}
			}
			return nil, nil
		},
		Description: "Returns the passed in argument converted into a string.",
	}}
}
