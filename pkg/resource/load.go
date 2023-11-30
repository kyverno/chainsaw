package resource

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	extyaml "github.com/kyverno/kyverno/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type (
	splitter  = func([]byte) ([][]byte, error)
	converter = func([]byte) ([]byte, error)
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

func LoadFromURL(url *url.URL) ([]unstructured.Unstructured, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tests, err := Parse(content)
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("found no test in %s", url.String())
	}
	return tests, nil
}

func Parse(content []byte) ([]unstructured.Unstructured, error) {
	return parse(content, nil, nil)
}

func parse(content []byte, splitter splitter, converter converter) ([]unstructured.Unstructured, error) {
	if splitter == nil {
		splitter = extyaml.SplitDocuments
	}
	if converter == nil {
		converter = yaml.ToJSON
	}

	documents, err := splitter(content)
	if err != nil {
		return nil, err
	}
	var resources []unstructured.Unstructured
	for _, document := range documents {
		jsonBytes, err := converter(document)
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
