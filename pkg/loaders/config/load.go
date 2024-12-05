package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/conversion"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/loaders"
	"github.com/kyverno/chainsaw/pkg/utils/fs"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	DefaultFileName = ".chainsaw.yaml"
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func() (loader.Loader, error)
	converter     = func(schema.GroupVersionKind, unstructured.Unstructured) (*v1alpha2.Configuration, error)
)

var (
	configuration_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("Configuration")
	configuration_v1alpha2 = v1alpha2.SchemeGroupVersion.WithKind("Configuration")
)

func Load(getter fs.Getter, path string) (*v1alpha2.Configuration, error) {
	path, err := getter.GetFile(path)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	config, err := LoadBytes(content)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration (%w)", err)
	}
	return config, nil
}

func LoadBytes(content []byte) (*v1alpha2.Configuration, error) {
	configs, err := Parse(content)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, errors.New("found no configuration")
	}
	if len(configs) != 1 {
		return nil, fmt.Errorf("found multiple configurations (%d)", len(configs))
	}
	return configs[0], nil
}

func Parse(content []byte) ([]*v1alpha2.Configuration, error) {
	return parse(content, nil, nil, nil)
}

func parse(content []byte, splitter splitter, loaderFactory loaderFactory, converter converter) ([]*v1alpha2.Configuration, error) {
	if splitter == nil {
		splitter = yaml.SplitDocuments
	}
	if loaderFactory == nil {
		loaderFactory = loaders.DefaultLoader
	}
	if converter == nil {
		converter = defaultConverter
	}
	documents, err := splitter(content)
	if err != nil {
		return nil, err
	}
	loader, err := loaderFactory()
	if err != nil {
		return nil, err
	}
	var configs []*v1alpha2.Configuration
	for _, document := range documents {
		gvk, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		config, err := converter(gvk, untyped)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func defaultConverter(gvk schema.GroupVersionKind, untyped unstructured.Unstructured) (*v1alpha2.Configuration, error) {
	scheme := runtime.NewScheme()
	if err := v1alpha1.Install(scheme); err != nil {
		return nil, err
	}
	if err := v1alpha2.Install(scheme); err != nil {
		return nil, err
	}
	if err := conversion.RegisterConversions(scheme); err != nil {
		return nil, err
	}
	var out v1alpha2.Configuration
	switch gvk {
	case configuration_v1alpha1:
		if err := scheme.Convert(&untyped, &out, nil); err != nil {
			return nil, err
		}
	case configuration_v1alpha2:
		if err := scheme.Convert(&untyped, &out, nil); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("type not supported %s", gvk)
	}
	return &out, nil
}
