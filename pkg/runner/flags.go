package runner

import (
	"flag"
	"fmt"
	"strconv"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func setupFlags(config v1alpha1.ConfigurationSpec) error {
	testing.Init()
	if err := flag.Set("test.v", "true"); err != nil {
		return err
	}
	if err := flag.Set("test.parallel", strconv.Itoa(config.Parallel)); err != nil {
		return err
	}
	if err := flag.Set("test.failfast", fmt.Sprint(config.FailFast)); err != nil {
		return err
	}
	if err := flag.Set("test.paniconexit0", "true"); err != nil {
		return err
	}
	if err := flag.Set("test.fullpath", "false"); err != nil {
		return err
	}
	if err := flag.Set("test.count", strconv.Itoa(*config.RepeatCount)); err != nil {
		return err
	}
	if err := flag.Set("test.run", config.IncludeTestRegex); err != nil {
		return err
	}
	if err := flag.Set("test.skip", config.ExcludeTestRegex); err != nil {
		return err
	}
	flag.Parse()
	return nil
}
