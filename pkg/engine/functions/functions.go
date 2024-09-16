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

type FunctionEntry struct {
	functions.FunctionEntry
	Note string
}

func GetFunctions() []FunctionEntry {
	return []FunctionEntry{{
		FunctionEntry: functions.FunctionEntry{
			Name: env,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpString}},
			},
			Handler: jpEnv,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: k8sGet,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpAny}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
			},
			Handler: jpKubernetesGet,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: k8sList,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpAny}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}, Optional: true},
			},
			Handler: jpKubernetesList,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: k8sExists,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpAny}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
			},
			Handler: jpKubernetesExists,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: k8sResourceExists,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpAny}},
				{Types: []functions.JpType{functions.JpString}},
				{Types: []functions.JpType{functions.JpString}},
			},
			Handler: jpKubernetesResourceExists,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: k8sServerVersion,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpAny}},
			},
			Handler: jpKubernetesServerVersion,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: metricsDecode,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpString}},
			},
			Handler: jpMetricsDecode,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
			Name: trimSpace,
			Arguments: []functions.ArgSpec{
				{Types: []functions.JpType{functions.JpString}},
			},
			Handler: jpTrimSpace,
		},
	}, {
		FunctionEntry: functions.FunctionEntry{
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
		},
	}}
}
