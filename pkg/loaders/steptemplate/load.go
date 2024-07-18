package steptemplate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	internalloader "github.com/kyverno/chainsaw/pkg/internal/loader"
	yamlutils "github.com/kyverno/chainsaw/pkg/utils/yaml"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/openapi"
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func(openapi.Client) (loader.Loader, error)
	converter     = func(unstructured.Unstructured) (*v1alpha1.StepTemplate, error)
)

func Load(path string, remarshal bool) ([]*v1alpha1.StepTemplate, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	tests, err := Parse(content, remarshal)
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("found no step template in %s", path)
	}
	return tests, nil
}

func Parse(content []byte, remarshal bool) ([]*v1alpha1.StepTemplate, error) {
	return parse(content, remarshal, nil, nil, nil)
}

func parse(content []byte, remarshal bool, splitter splitter, loaderFactory loaderFactory, converter converter) ([]*v1alpha1.StepTemplate, error) {
	if splitter == nil {
		splitter = yaml.SplitDocuments
	}
	if converter == nil {
		converter = convert.To[v1alpha1.StepTemplate]
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
	var tests []*v1alpha1.StepTemplate
	for _, document := range documents {
		if remarshal {
			altered, err := yamlutils.Remarshal(document)
			if err != nil {
				return nil, err
			}
			document = altered
		}
		_, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		test, err := converter(untyped)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	return tests, nil
}
