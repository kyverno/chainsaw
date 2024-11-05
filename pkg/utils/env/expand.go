package env

import (
	"os"
)

func Expand(env map[string]string, in ...string) []string {
	var args []string
	for _, arg := range in {
		expanded := os.Expand(arg, func(key string) string {
			expanded := env[key]
			if key == "$" {
				expanded = "$"
			} else if expanded == "" {
				expanded = os.Getenv(key)
			}
			return expanded
		})
		if expanded != "" {
			args = append(args, expanded)
		} else {
			args = append(args, arg)
		}
	}
	return args
}
