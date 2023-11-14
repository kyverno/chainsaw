package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno/ext/output/color"
)

type operation struct {
	command   v1alpha1.Command
	namespace string
	log       bool
}

func New(command v1alpha1.Command, namespace string, log bool) *operation {
	return &operation{
		command:   command,
		namespace: namespace,
		log:       log,
	}
}

func (c *operation) Cleanup() {}

func (c *operation) Exec(ctx context.Context) (_err error) {

	logger := logging.FromContext(ctx)
	const operation = "CMD   "
	var output internal.CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(operation, color.BoldGreen, "DONE")
		} else {
			logger.Log(operation, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	if c.log {
		defer func() {
			if out := output.Out(); out != "" {
				logger.Log("STDOUT", color.BoldFgCyan, "LOGS...\n"+out)
			}
			if err := output.Err(); err != "" {
				logger.Log("STDERR", color.BoldFgCyan, "LOGS...\n"+err)
			}
		}()
	} else {
		logger.Log("STD___", color.BoldYellow, "suppressed logs")
	}
	args := expand(map[string]string{"NAMESPACE": c.namespace}, c.command.Args...)
	cmd := exec.CommandContext(ctx, c.command.Entrypoint, args...) //nolint:gosec
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAMESPACE=%s", c.namespace))
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	logger.Log(operation, color.BoldFgCyan, cmd, "RUNNING...")
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	return cmd.Run()
}

func expand(env map[string]string, in ...string) []string {
	var args []string
	for _, arg := range in {
		expanded := os.Expand(arg, func(key string) string {
			expanded := env[key]
			if expanded == "" {
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
