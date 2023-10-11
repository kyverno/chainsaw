package flag

import (
	"github.com/spf13/pflag"
)

// IsSet returns true if a flag is set on the command line.
func IsSet(flagSet *pflag.FlagSet, name string) bool {
	found := false
	flagSet.Visit(func(flag *pflag.Flag) {
		if flag.Name == name {
			found = true
		}
	})
	return found
}
