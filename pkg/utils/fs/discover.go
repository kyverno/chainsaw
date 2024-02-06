package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"
)

func DiscoverFolders(paths ...string) ([]string, error) {
	folders := sets.New[string]()
	var errors []error

	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil {
			errors = append(errors, fmt.Errorf("error checking path %s: %v", path, err))
			continue
		}
		err = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				folders.Insert(file)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return sets.List(folders), nil
}
