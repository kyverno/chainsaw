package discovery

import (
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
)

func DiscoverTests(fileName string, paths ...string) ([]Test, error) {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		return nil, err
	}
	var tests []Test
	for _, folder := range folders {
		test, err := LoadTest(fileName, folder)
		if err != nil {
			return nil, err
		} else if test != nil {
			tests = append(tests, *test)
		}
	}
	return tests, nil
}
