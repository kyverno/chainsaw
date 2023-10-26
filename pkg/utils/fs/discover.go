package fs

import (
	"io/fs"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/util/sets"
)

func DiscoverFolders(paths ...string) ([]string, error) {
	folders := sets.New[string]()
	for _, path := range paths {
		if _, err := os.Lstat(path); err == nil {
			err := filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
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
	}
	return sets.List(folders), nil
}
