package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	internalloader "github.com/kyverno/chainsaw/pkg/internal/loader"
	testvalidation "github.com/kyverno/chainsaw/pkg/validation/test"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	yamlv3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/openapi"
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func(openapi.Client) (loader.Loader, error)
	converter     = func(unstructured.Unstructured) (*v1alpha1.Test, error)
	validator     = func(obj *v1alpha1.Test) field.ErrorList
)

var test_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("Test")

func Load(path string) ([]*v1alpha1.Test, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	tests, err := Parse(content)
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("found no test in %s", path)
	}
	return tests, nil
}

func Parse(content []byte) ([]*v1alpha1.Test, error) {
	return parse(content, nil, nil, nil, nil)
}

func remarshal(document []byte) ([]byte, error) {
	var pre map[string]any
	err := yamlv3.Unmarshal(document, &pre)
	if err != nil {
		return nil, err
	}
	return yamlv3.Marshal(pre)
}

func parse(content []byte, splitter splitter, loaderFactory loaderFactory, converter converter, validator validator) ([]*v1alpha1.Test, error) {
	if splitter == nil {
		splitter = yaml.SplitDocuments
	}
	if converter == nil {
		converter = convert.To[v1alpha1.Test]
	}
	if validator == nil {
		validator = testvalidation.ValidateTest
	}
	var loader loader.Loader
	if loaderFactory != nil {
		_loader, err := loaderFactory(internalloader.OpenApiClient)
		if err != nil {
			return nil, err
		}
		loader = _loader
	}
	if loader == nil {
		if internalloader.Err != nil {
			return nil, internalloader.Err
		}
		loader = internalloader.DefaultLoader
	}
	documents, err := splitter(content)
	if err != nil {
		return nil, err
	}
	var tests []*v1alpha1.Test
	for _, document := range documents {
		document, err := remarshal(document)
		if err != nil {
			return nil, err
		}
		gvk, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		switch gvk {
		case test_v1alpha1:
			test, err := converter(untyped)
			if err != nil {
				return nil, err
			}
			if err := validator(test).ToAggregate(); err != nil {
				return nil, err
			}
			tests = append(tests, test)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return tests, nil
}
