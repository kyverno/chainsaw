package resource

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-getter"
	extyaml "github.com/kyverno/kyverno/ext/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type (
	splitter  = func([]byte) ([][]byte, error)
	converter = func([]byte) ([]byte, error)
)

func Load(path string, manifest bool) ([]unstructured.Unstructured, error) {
	var resources []unstructured.Unstructured
	matchingFiles, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("failed to match files \"%s\": %v", path, err)
	}
	if len(matchingFiles) == 0 {
		return nil, fmt.Errorf("no files found matching path: %s", path)
	}
	for _, file := range matchingFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		tests, err := Parse(content, manifest)
		if err != nil {
			return nil, err
		}
		if len(tests) == 0 {
			return nil, fmt.Errorf("found no resource in %s", file)
		}
		resources = append(resources, tests...)
	}
	return resources, nil
}

func LoadFromURI(url *url.URL, manifest bool) ([]unstructured.Unstructured, error) {
	tempFile, err := os.CreateTemp("", "getter-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %s", err)
	}
	defer os.Remove(tempFile.Name())
	client := &getter.Client{
		Ctx:  context.Background(),
		Src:  url.String(),
		Dst:  tempFile.Name(),
		Mode: getter.ClientModeFile,
	}
	backoff := wait.Backoff{
		Steps:    3,
		Duration: 1 * time.Second,
		Factor:   2.0,
		Jitter:   0.1,
	}
	if err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		if err := client.Get(); err != nil {
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("error downloading content: %s", err)
	}
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading downloaded content: %s", err)
	}
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("error closing temp file: %s", err)
	}
	tests, err := Parse(content, manifest)
	if err != nil {
		return nil, err
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("found no test in %s", url.String())
	}
	return tests, nil
}

func Parse(content []byte, manifest bool) ([]unstructured.Unstructured, error) {
	return parse(content, nil, nil, manifest)
}

func parse(content []byte, splitter splitter, converter converter, manifest bool) ([]unstructured.Unstructured, error) {
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
			if manifest || !runtime.IsMissingKind(err) {
				return nil, err
			}
		}
		resources = append(resources, resource)
	}
	return resources, nil
}
