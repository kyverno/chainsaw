package fs

import (
	"errors"
	"os"
)

func CheckFolders(paths ...string) error {
	if len(paths) == 0 {
		return errors.New("no folders provided")
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}
