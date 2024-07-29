package steptemplate

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
)

type (
	splitter      = func([]byte) ([][]byte, error)
	loaderFactory = func() (loader.Loader, error)
	converter     = func(unstructured.Unstructured) (*v1alpha1.StepTemplate, error)
)

var stepTemplate_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("StepTemplate")

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
	if loaderFactory == nil {
		loaderFactory = loaders.DefaultLoader
	}
	if converter == nil {
		converter = convert.To[v1alpha1.StepTemplate]
	}
	documents, err := splitter(content)
	if err != nil {
		return nil, err
	}
	loader, err := loaderFactory()
	if err != nil {
		return nil, err
	}
	var steps []*v1alpha1.StepTemplate
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
		case stepTemplate_v1alpha1:
			step, err := converter(untyped)
			if err != nil {
				return nil, err
			}
			steps = append(steps, step)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return steps, nil
}
