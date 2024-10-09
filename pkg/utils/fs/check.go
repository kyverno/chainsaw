package fs

import (
	"net/url"
	"os"
)

func CheckFolders(paths ...string) error {
	for _, path := range paths {
		base, err := url.Parse(path)
		if err != nil {
			return err
		}
		if base.Scheme != "" {
			continue
		}
		if _, err := os.Stat(base.Path); err != nil {
			return err
		}
	}
	return nil
}
