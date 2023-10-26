package v1alpha1

// Command describes a command to run as a part of a test step or suite.
type Command struct {
	// The command and argument to run as a string.
	Command string `json:"command,omitempty"`
	// If set, the `--namespace` flag will be appended to the command with the namespace to use.
	Namespaced bool `json:"namespaced,omitempty"`
	// Ability to run a shell script from TestStep (without a script file)
	// namespaced and command should not be used with script.  namespaced is ignored and command is an error.
	// env expansion is depended upon the shell but ENV is passed to the runtime env.
	Script string `json:"script,omitempty"`
	// If set, exit failures (`exec.ExitError`) will be ignored. `exec.Error` are NOT ignored.
	ContinueOnError *bool `json:"continueOnError,omitempty"`
	// Override the TestSuite timeout for this command (in seconds).
	Timeout int `json:"timeout,omitempty"`
	// If set, the output from the command is NOT logged.  Useful for sensitive logs or to reduce noise.
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`
}
