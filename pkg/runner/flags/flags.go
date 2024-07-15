package flags

import (
	"strconv"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
)

func GetFlags(config v1alpha2.ConfigurationSpec) map[string]string {
	flags := map[string]string{
		"test.v":            "true",
		"test.paniconexit0": "true",
		"test.fullpath":     "false",
		"test.run":          config.Discovery.IncludeTestRegex,
		"test.skip":         config.Discovery.ExcludeTestRegex,
	}
	if config.Execution.Parallel != nil {
		flags["test.parallel"] = strconv.Itoa(*config.Execution.Parallel)
	}
	if config.Execution.RepeatCount != nil {
		flags["test.count"] = strconv.Itoa(*config.Execution.RepeatCount)
	}
	return flags
}
