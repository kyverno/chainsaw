package discovery

import (
	fsutils "github.com/kyverno/chainsaw/pkg/utils/fs"
)

type folders = func() []string

func DiscoverTests(fileName string, paths ...string) ([]Test, error) {
	folders, err := fsutils.DiscoverFolders(paths...)
	if err != nil {
		return nil, err
	}
	return helpDiscoverTests(fileName, func() []string {
		return folders
	})
}

func helpDiscoverTests(fileName string, discoveredFolders folders) ([]Test, error) {
	var tests []Test
	for _, folder := range discoveredFolders() {
		test, err := LoadTest(fileName, folder)
		if err != nil {
			return nil, err
		} else if test != nil {
			tests = append(tests, *test)
		}
	}
	return tests, nil
}
