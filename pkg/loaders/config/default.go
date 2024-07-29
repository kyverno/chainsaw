package config

import (
	"sync"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/data"
)

func defaultConfiguration(_fs func() ([]byte, error)) (*v1alpha2.Configuration, error) {
	if _fs == nil {
		_fs = data.ConfigFile
	}
	bytes, err := _fs()
	if err != nil {
		return nil, err
	}
	return LoadBytes(bytes)
}

var DefaultConfiguration = sync.OnceValues(func() (*v1alpha2.Configuration, error) { return defaultConfiguration(nil) })
