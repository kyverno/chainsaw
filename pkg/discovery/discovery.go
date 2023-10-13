package discovery

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/test"
)

type Test struct {
	*v1alpha1.Test
	BasePath string
}

func DiscoverTests(fileName string, paths ...string) ([]Test, error) {
	var files []string
	for _, path := range paths {
		if _, err := os.Lstat(path); err == nil {
			err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.Name() == fileName {
					files = append(files, file)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	var tests []Test
	for _, file := range files {
		apiTests, err := test.Load(file)
		if err == nil {
			basePath := filepath.Dir(file)
			for _, test := range apiTests {
				tests = append(tests, Test{
					Test:     test,
					BasePath: basePath,
				})
			}
		}
	}
	return tests, nil
}
