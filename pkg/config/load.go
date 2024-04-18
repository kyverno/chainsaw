package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/data"
	configvalidation "github.com/kyverno/chainsaw/pkg/validation/config"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

const (
	DefaultFileName = ".chainsaw.yaml"
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func(openapi.Client) (loader.Loader, error)
	converter     = func(unstructured.Unstructured) (*v1alpha1.Configuration, error)
	validator     = func(obj *v1alpha1.Configuration) field.ErrorList
)

var configuration_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("Configuration")

func Load(path string) (*v1alpha1.Configuration, error) {
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

func LoadBytes(content []byte) (*v1alpha1.Configuration, error) {
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

func Parse(content []byte) ([]*v1alpha1.Configuration, error) {
	return parse(content, nil, nil, nil, nil)
}

func parse(content []byte, splitter splitter, loaderFactory loaderFactory, converter converter, validator validator) ([]*v1alpha1.Configuration, error) {
	if splitter == nil {
		splitter = yaml.SplitDocuments
	}
	if loaderFactory == nil {
		loaderFactory = loader.New
	}
	if converter == nil {
		converter = convert.To[v1alpha1.Configuration]
	}
	if validator == nil {
		validator = configvalidation.ValidateConfiguration
	}
	documents, err := splitter(content)
	if err != nil {
		return nil, err
	}
	var configs []*v1alpha1.Configuration
	// TODO: no need to allocate a validator every time
	loader, err := loaderFactory(openapiclient.NewLocalCRDFiles(data.Crds(), data.CrdsFolder))
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		gvk, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		switch gvk {
		case configuration_v1alpha1:
			config, err := converter(untyped)
			if err != nil {
				return nil, err
			}
			if err := validator(config).ToAggregate(); err != nil {
				return nil, err
			}
			configs = append(configs, config)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return configs, nil
}
