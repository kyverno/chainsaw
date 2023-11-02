package flags

import (
	"fmt"
	"strconv"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func GetFlags(config v1alpha1.ConfigurationSpec) map[string]string {
	flags := map[string]string{
		"test.v":            "true",
		"test.paniconexit0": "true",
		"test.fullpath":     "false",
		"test.failfast":     fmt.Sprint(config.FailFast),
		"test.run":          config.IncludeTestRegex,
		"test.skip":         config.ExcludeTestRegex,
	}
	if config.Parallel != nil {
		flags["test.parallel"] = strconv.Itoa(*config.Parallel)
	}
	if config.RepeatCount != nil {
		flags["test.count"] = strconv.Itoa(*config.RepeatCount)
	}
	return flags
}
