package validation

import (
	"net/url"
	"os"
)

func IsValidPathOrURI(pathStr string) bool {
	if _, err := os.Stat(pathStr); err == nil {
		return true
	}
	if _, err := url.ParseRequestURI(pathStr); err == nil {
		return true
	}
	return false
}
