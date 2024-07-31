package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"
)

func discoverFolders(stat func(string) (os.FileInfo, error), walk func(string, filepath.WalkFunc) error, paths ...string) ([]string, error) {
	if stat == nil {
		stat = os.Stat
	}
	if walk == nil {
		walk = filepath.Walk
	}
	folders := sets.New[string]()
	var errors []error
	for _, path := range paths {
		_, err := stat(path)
		if err != nil {
			errors = append(errors, fmt.Errorf("error checking path %s: %v", path, err))
			continue
		}
		err = walk(path, func(file string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				folders.Insert(file)
			}
			return nil
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return sets.List(folders), nil
}

func DiscoverFolders(paths ...string) ([]string, error) {
	return discoverFolders(nil, nil, paths...)
}
