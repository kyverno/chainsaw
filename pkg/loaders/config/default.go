package config

import (
	"sync"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/data"
)

func defaultConfiguration() (*v1alpha2.Configuration, error) {
	bytes, err := data.ConfigFile()
	if err != nil {
		return nil, err
	}
	return LoadBytes(bytes)
}

var DefaultConfiguration = sync.OnceValues(defaultConfiguration)
