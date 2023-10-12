package resource

import (
	"fmt"
	"os"
	"path/filepath"

	yamlutils "github.com/kyverno/chainsaw/pkg/utils/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Load(path string) ([]unstructured.Unstructured, error) {
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

func Parse(content []byte) ([]unstructured.Unstructured, error) {
	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var resources []unstructured.Unstructured
	for _, document := range documents {
		jsonBytes, err := yaml.ToJSON(document)
		if err != nil {
			return nil, err
		}
		var resource unstructured.Unstructured
		if err := resource.UnmarshalJSON(jsonBytes); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}
