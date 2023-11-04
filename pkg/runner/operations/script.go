package operations

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
)

func operationScript(ctx context.Context, script v1alpha1.Script, log bool, namespace string) (_err error) {
	logger := logging.FromContext(ctx)
	const operation = "SCRIPT"
	var output CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(operation, color.BoldGreen, "DONE")
		} else {
			logger.Log(operation, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	if log {
		defer func() {
			if out := output.Out(); out != "" {
				logger.Log("STDOUT", color.BoldFgCyan, "LOGS...\n"+out)
			}
			if err := output.Err(); err != "" {
				logger.Log("STDERR", color.BoldFgCyan, "LOGS...\n"+err)
			}
		}()
	} else {
		logger.Log("STDXXX", color.BoldYellow, "suppressed logs")
	}
	cmd := exec.CommandContext(ctx, "sh", "-c", script.Content) //nolint:gosec
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAMESPACE=%s", namespace))
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	logger.Log(operation, color.BoldFgCyan, "RUNNING...")
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr
	return cmd.Run()
}
