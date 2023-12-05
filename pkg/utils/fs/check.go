package fs

import (
	"os"
)

func CheckFolders(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}
