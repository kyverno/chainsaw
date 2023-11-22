package step

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	internalloader "github.com/kyverno/chainsaw/pkg/internal/loader"
	"github.com/kyverno/chainsaw/pkg/validation"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/kyverno/kyverno/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/openapi"
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func(openapi.Client) (loader.Loader, error)
	converter     = func(unstructured.Unstructured) (*v1alpha1.TestStep, error)
)

var testStep_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("TestStep")

func Load(path string) ([]*v1alpha1.TestStep, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	testSteps, err := Parse(content)
	if err != nil {
		return nil, err
	}
	if len(testSteps) == 0 {
		return nil, fmt.Errorf("found no testStep in %s", path)
	}
	return testSteps, nil
}

func Parse(content []byte) ([]*v1alpha1.TestStep, error) {
	return parse(content, nil, nil, nil)
}

func parse(content []byte, splitter splitter, loaderFactory loaderFactory, converter converter) ([]*v1alpha1.TestStep, error) {
	if splitter == nil {
		splitter = yaml.SplitDocuments
	}
	if converter == nil {
		converter = convert.To[v1alpha1.TestStep]
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
	var testSteps []*v1alpha1.TestStep
	for _, document := range documents {
		gvk, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		switch gvk {
		case testStep_v1alpha1:
			testStep, err := converter(untyped)
			if err != nil {
				return nil, err
			}
			if err := validation.ValidateTestStep(testStep).ToAggregate(); err != nil {
				return nil, err
			}
			testSteps = append(testSteps, testStep)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return testSteps, nil
}
