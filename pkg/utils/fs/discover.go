package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

func DiscoverFolders(paths ...string) ([]string, error) {
	folders := sets.New[string]()
	var errorMessages []string

	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("error checking path %s: %v", path, err))
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
			errorMessages = append(errorMessages, fmt.Sprintf("error walking the path %s: %v", path, err))
		}
	}
	if len(errorMessages) > 0 {
		return nil, errors.New(strings.Join(errorMessages, "; "))
	}
	return sets.List(folders), nil
}
