package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/kyverno/kyverno/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

const (
	DefaultFileName = ".chainsaw.yaml"
)

type (
	DocumentParser = func([]byte) ([][]byte, error)
	crdFolder      = func(string) string
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
	yamlDocumentParser := func(content []byte) ([][]byte, error) {
		return yaml.SplitDocuments(content)
	}
	configs, err := Parse(content, yamlDocumentParser)
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

func Parse(content []byte, yamlDocumentParser DocumentParser) ([]*v1alpha1.Configuration, error) {
	documents, err := yamlDocumentParser(content)
	if err != nil {
		return nil, err
	}
	crdFolder := func(_ string) string {
		return data.CrdsFolder
	}
	return parseDocuments(documents, crdFolder)
}

func parseDocuments(documents [][]byte, crdfolder crdFolder) ([]*v1alpha1.Configuration, error) {
	var configurations []*v1alpha1.Configuration
	loader, err := newLoader(data.Crds(), crdfolder(""))
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		configuration, err := parseDocument(loader, document)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, configuration)
	}
	return configurations, nil
}

func newLoader(crds fs.FS, crdsFolder string) (loader.Loader, error) {
	return loader.New(openapiclient.NewLocalCRDFiles(crds, crdsFolder))
}

func parseDocument(loader loader.Loader, document []byte) (*v1alpha1.Configuration, error) {
	gvk, untyped, err := loader.Load(document)
	if err != nil {
		return nil, err
	}
	if gvk != configuration_v1alpha1 {
		return nil, fmt.Errorf("type not supported %s", gvk)
	}
	return convert.To[v1alpha1.Configuration](untyped)
}
