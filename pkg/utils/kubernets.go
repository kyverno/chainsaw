package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// LoadYAML loads all objects from a reader
func LoadYAML(path string, r io.Reader) ([]client.Object, error) {
	yamlReader := yaml.NewYAMLReader(bufio.NewReader(r))

	objects := []client.Object{}

	for {
		data, err := yamlReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading yaml %s: %w", path, err)
		}

		unstructuredObj := &unstructured.Unstructured{}
		decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(data), len(data))

		if err = decoder.Decode(unstructuredObj); err != nil {
			return nil, fmt.Errorf("error decoding yaml %s: %w", path, err)
		}
		var obj client.Object

		err = runtime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(unstructuredObj.Object, obj, true)
		if err != nil {
			return nil, fmt.Errorf("error converting unstructured object %s %v", path, err)
		}
		// discovered reader will return empty objects if a number of lines are preceding a yaml separator (---)
		// this detects that, logs and continues
		if obj.GetObjectKind().GroupVersionKind().Kind == "" {
			log.Println("object detected with no GVK Kind for path", path)
		} else {
			objects = append(objects, obj)
		}
	}

	return objects, nil
}

// LoadYAMLFromFile loads all objects from a YAML file.
func LoadYAMLFromFile(path string) ([]client.Object, error) {
	opened, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer opened.Close()

	return LoadYAML(path, opened)
}
