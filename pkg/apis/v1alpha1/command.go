package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Command describes a command to run as a part of a test step or suite.
type Command struct {
	// The command and argument to run as a string.
	// +optional
	Command string `json:"command,omitempty"`
	// If set, the `--namespace` flag will be appended to the command with the namespace to use.
	// +optional
	Namespaced bool `json:"namespaced,omitempty"`
	// Ability to run a shell script from TestStep (without a script file)
	// namespaced and command should not be used with script.  namespaced is ignored and command is an error.
	// +optional
	Script string `json:"script,omitempty"`
	// If set, exit failures (`exec.ExitError`) will be ignored. `exec.Error` are NOT ignored.
	// +optional
	ContinueOnError *bool `json:"continueOnError,omitempty"`
	// Override the TestSuite timeout for this command (in seconds).
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// If set, the output from the command is NOT logged.  Useful for sensitive logs or to reduce noise.
	// +optional
	SkipLogOutput bool `json:"skipLogOutput,omitempty"`
}
