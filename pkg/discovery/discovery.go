package discovery

import (
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
	"k8s.io/apimachinery/pkg/labels"
)

func DiscoverTests(fileName string, selector labels.Selector, remarshal bool, paths ...string) ([]Test, error) {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		return nil, err
	}
	return discoverTests(fileName, selector, remarshal, folders...)
}

func discoverTests(fileName string, selector labels.Selector, remarshal bool, folders ...string) ([]Test, error) {
	if selector == nil {
		selector = labels.Everything()
	}
	var tests []Test
	for _, folder := range folders {
		t, err := LoadTest(fileName, folder, remarshal)
		if err != nil {
			return nil, err
		}
		for _, t := range t {
			if selector.Matches(labels.Set(t.Test.Labels)) {
				tests = append(tests, t)
			}
		}
	}
	return tests, nil
}
