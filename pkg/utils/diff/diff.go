package diff

import (
	"github.com/pmezard/go-difflib/difflib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func prune(expected map[string]any, actual map[string]any) {
	for k, v := range actual {
		if expected[k] == nil {
			delete(actual, k)
		} else {
			if v, ok := v.(map[string]any); ok {
				prune(expected[k].(map[string]any), v)
			}
		}
	}
}

func pruneMetadata(expected map[string]any, actual map[string]any) {
	for k, v := range actual {
		if k == "ownerReferences" || k == "name" {
			continue
		}
		if expected[k] == nil {
			delete(actual, k)
		} else {
			if v, ok := v.(map[string]any); ok {
				prune(expected[k].(map[string]any), v)
			}
		}
	}
}

func pruneRoot(expected map[string]any, actual map[string]any) {
	for k, v := range actual {
		if expected[k] == nil {
			delete(actual, k)
		} else {
			if v, ok := v.(map[string]any); ok {
				if k == "metadata" {
					pruneMetadata(expected[k].(map[string]any), v)
				} else {
					prune(expected[k].(map[string]any), v)
				}
			}
		}
	}
}

func PrettyDiff(expected unstructured.Unstructured, actual unstructured.Unstructured) (string, error) {
	pruneRoot(expected.Object, actual.Object)
	if expectedBuf, err := yaml.Marshal(&expected); err != nil {
		return "", err
	} else if actualBuf, err := yaml.Marshal(&actual); err != nil {
		return "", err
	} else {
		diffed := difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(expectedBuf)),
			B:        difflib.SplitLines(string(actualBuf)),
			FromFile: "expected",
			ToFile:   "actual",
			Context:  3,
		}
		return difflib.GetUnifiedDiffString(diffed)
	}
}
