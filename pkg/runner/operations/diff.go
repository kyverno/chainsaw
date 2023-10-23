package operations

import (
	"fmt"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func diff(expected, actual unstructured.Unstructured) (string, error) {
	expectedBytes, err := yaml.Marshal(expected)
	if err != nil {
		return "", fmt.Errorf("failed to marshal expected content to YAML: %w", err)
	}
	candidateBytes, err := yaml.Marshal(actual)
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
	diffStr, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", fmt.Errorf("failed to generate unified diff string: %w", err)
	}
	return strings.TrimSpace(diffStr), nil
}
