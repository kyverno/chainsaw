package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/kyverno/kyverno/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var test_v1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("Test")

func Load(path string) ([]*v1alpha1.Test, error) {
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

func Parse(content []byte) ([]*v1alpha1.Test, error) {
	documents, err := yaml.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var tests []*v1alpha1.Test
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
		case test_v1alpha1:
			test, err := convert.To[v1alpha1.Test](untyped)
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
