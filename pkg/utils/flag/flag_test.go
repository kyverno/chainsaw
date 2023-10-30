package flag

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestIsSet(t *testing.T) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flagSet.String("test-flag", "", "a test flag")
	flagSet.String("another-flag", "", "another test flag")
	flagSet.Bool("bool-flag", false, "a boolean flag")

	tests := []struct {
		name           string
		args           []string
		flagToCheck    string
		expectedResult bool
	}{
		{
			name:           "Single flag set",
			args:           []string{"--test-flag=value"},
			flagToCheck:    "test-flag",
			expectedResult: true,
		},
		{
			name:           "Multiple flags set - check unset flag",
			args:           []string{"--test-flag=value", "--bool-flag"},
			flagToCheck:    "another-flag",
			expectedResult: false,
		},
		{
			name:           "Multiple flags set - check set flag",
			args:           []string{"--test-flag=value", "--bool-flag"},
			flagToCheck:    "bool-flag",
			expectedResult: true,
		},
		{
			name:           "No flags set",
			args:           []string{},
			flagToCheck:    "test-flag",
			expectedResult: false,
		},
		{
			name:           "Check undefined flag",
			args:           []string{},
			flagToCheck:    "undefined-flag",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localFlagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			localFlagSet.String("test-flag", "", "a test flag")
			localFlagSet.String("another-flag", "", "another test flag")
			localFlagSet.Bool("bool-flag", false, "a boolean flag")

			err := localFlagSet.Parse(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, IsSet(localFlagSet, tt.flagToCheck))
		})
	}
}
