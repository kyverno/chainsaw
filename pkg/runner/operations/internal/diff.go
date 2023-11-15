package internal

import (
	"fmt"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

type (
	diffLibInterface = func(difflib.UnifiedDiff) (string, error)
	yamlMarshaler    func(obj interface{}) ([]byte, error)
)

func diff(expected, actual ctrlclient.Object) (string, error) {
	return diffHelper(expected, actual, nil, nil, nil)
}

func diffHelper(expected, actual interface{}, expectedMarshaler, actualMarshaler yamlMarshaler, getDiffString diffLibInterface) (string, error) {
	if expectedMarshaler == nil {
		expectedMarshaler = yaml.Marshal
	}
	if actualMarshaler == nil {
		actualMarshaler = yaml.Marshal
	}
	if getDiffString == nil {
		getDiffString = difflib.GetUnifiedDiffString
	}

	expectedBytes, err := expectedMarshaler(expected)
	if err != nil {
		return "", fmt.Errorf("failed to marshal expected content to YAML: %w", err)
	}
	candidateBytes, err := actualMarshaler(actual)
	if err != nil {
		return "", fmt.Errorf("failed to marshal candidate content to YAML: %w", err)
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(expectedBytes)),
		B:        difflib.SplitLines(string(candidateBytes)),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  3,
	}

	diffStr, err := getDiffString(diff)
	if err != nil {
		return "", fmt.Errorf("failed to generate unified diff string: %w", err)
	}
	return strings.TrimSpace(diffStr), nil
}
