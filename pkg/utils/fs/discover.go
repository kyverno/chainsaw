package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"
)

func discoverFolders(getter Getter, stat func(string) (os.FileInfo, error), walk func(string, filepath.WalkFunc) error, paths ...string) ([]string, error) {
	if stat == nil {
		stat = os.Stat
	}
	if walk == nil {
		walk = filepath.Walk
	}
	folders := sets.New[string]()
	var errors []error
	for _, path := range paths {
		path, err := getter.Get(path)
		if err != nil {
			return nil, err
		}
		if _, err := stat(path); err != nil {
			errors = append(errors, fmt.Errorf("error checking path %s: %v", path, err))
			continue
		}
		if err := walk(path, func(file string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				folders.Insert(file)
			}
			return nil
		}); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return nil, multierr.Combine(errors...)
	}
	return sets.List(folders), nil
}

func DiscoverFolders(getter Getter, paths ...string) ([]string, error) {
	return discoverFolders(getter, nil, nil, paths...)
}
