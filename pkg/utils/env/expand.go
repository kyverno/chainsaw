package env

import (
	"os"
)

func Expand(env map[string]string, in ...string) []string {
	var args []string
	for _, arg := range in {
		expanded := os.Expand(arg, func(key string) string {
			// check $ key -> $$ becomes $
			if key == "$" {
				return "$"
			}
			// check the env map first
			if expanded, found := env[key]; found {
				return expanded
			}
			// check the env vars
			if expanded, found := os.LookupEnv(key); found {
				return expanded
			}
			// return the given key
			return "$" + key
		})
		if expanded != "" {
			args = append(args, expanded)
		} else {
			args = append(args, arg)
		}
	}
	return args
}
