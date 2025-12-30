package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/loaders"
	yamlutils "github.com/kyverno/chainsaw/pkg/utils/yaml"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func() (loader.Loader, error)
	converter     = func(unstructured.Unstructured) (*v1alpha1.Test, error)
)

var test_v1alpha1 = schema.GroupVersion(v1alpha1.GroupVersion).WithKind("Test")

func Load(file string, remarshal bool) ([]*v1alpha1.Test, error) {
	content, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	tests, err := Parse(content, remarshal)
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("found no test in %s", file)
	}
	return tests, nil
}

func Parse(content []byte, remarshal bool) ([]*v1alpha1.Test, error) {
	return parse(content, remarshal, nil, nil, nil)
}

func parse(content []byte, remarshal bool, splitter splitter, loaderFactory loaderFactory, converter converter) ([]*v1alpha1.Test, error) {
	if splitter == nil {
		splitter = yaml.SplitDocuments
	}
	if loaderFactory == nil {
		loaderFactory = loaders.DefaultLoader
	}
	if converter == nil {
		converter = convert.To[v1alpha1.Test]
	}
	documents, err := splitter(content)
	if err != nil {
		return nil, err
	}
	loader, err := loaderFactory()
	if err != nil {
		return nil, err
	}
	var tests []*v1alpha1.Test
	for _, document := range documents {
		if remarshal {
			altered, err := yamlutils.Remarshal(document)
			if err != nil {
				return nil, err
			}
			document = altered
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
			tests = append(tests, test)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return tests, nil
}
