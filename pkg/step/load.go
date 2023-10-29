package step

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/data"
	"github.com/kyverno/chainsaw/pkg/utils/loader"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
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
	documents, err := yaml.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var testSteps []*v1alpha1.TestStep
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
		case testStep_v1alpha1:
			testStep, err := convert.To[v1alpha1.TestStep](untyped)
			if err != nil {
				return nil, err
			}
			testSteps = append(testSteps, testStep)
		default:
			return nil, fmt.Errorf("type not supported %s", gvk)
		}
	}
	return testSteps, nil
}
