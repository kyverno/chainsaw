package flags

import (
	"flag"
	"strconv"
	"testing"

	"github.com/kyverno/chainsaw/pkg/model"
)

func getFlags(config model.Configuration) map[string]string {
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

func SetupFlags(config model.Configuration) error {
	testing.Init()
	for k, v := range getFlags(config) {
		if err := flag.Set(k, v); err != nil {
			return err
		}
	}
	flag.Parse()
	return nil
}
