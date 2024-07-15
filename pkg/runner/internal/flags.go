package internal

import (
	"flag"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/runner/flags"
)

func SetupFlags(config v1alpha2.ConfigurationSpec) error {
	testing.Init()
	for k, v := range flags.GetFlags(config) {
		if err := flag.Set(k, v); err != nil {
			return err
		}
	}
	flag.Parse()
	return nil
}
