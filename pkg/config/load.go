package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/chainsaw/pkg/utils/convert"
	"github.com/kyverno/chainsaw/pkg/utils/loader"
	yamlutils "github.com/kyverno/chainsaw/pkg/utils/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var configuration_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("Configuration")

func Load(path string) (*v1alpha1.Configuration, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	configs, err := Parse(content)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("Found no configuration in %s", path)
	}
	if len(configs) != 1 {
		return nil, fmt.Errorf("Found multiple configurations in %s (%d)", path, len(configs))
	}
	return configs[0], nil
}

func Parse(content []byte) ([]*v1alpha1.Configuration, error) {
	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var policies []*v1alpha1.Configuration
	// TODO: no need to allocate a validator every time
	loader, err := loader.New(openapiclient.NewLocalCRDFiles(data.Crds(), data.CrdsFolder))
	if err != nil {
		return nil, err
	}
	for _, document := range documents {
		gvk, untyped, err := loader.Load(document)
		if err != nil {
			return nil, err
		}
		switch gvk {
		case configuration_v1alpha1:
			policy, err := convert.To[v1alpha1.Configuration](untyped)
			if err != nil {
				return nil, err
			}
			policies = append(policies, policy)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return policies, nil
}
