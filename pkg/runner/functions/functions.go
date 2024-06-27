package functions

import (
	"github.com/jmespath-community/go-jmespath/pkg/functions"
)

var (
	// stable functions
	env       = stable("env")
	trimSpace = stable("trim_space")
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
		Handler: jpEnv,
	}, {
		Name: k8sGet,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpKubernetesGet,
	}, {
		Name: k8sList,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpKubernetesList,
	}, {
		Name: k8sExists,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpKubernetesExists,
	}, {
		Name: k8sResourceExists,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpKubernetesResourceExists,
	}, {
		Name: k8sServerVersion,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
		},
		Handler: jpKubernetesServerVersion,
	}, {
		Name: metricsDecode,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpMetricsDecode,
	}, {
		Name: trimSpace,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpTrimSpace,
	}}
}
