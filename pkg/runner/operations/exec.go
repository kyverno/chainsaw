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

func Exec(ctx context.Context, logger logging.Logger, exec v1alpha1.Exec, log bool, namespace string) error {
	if exec.Command != nil {
		return command(ctx, logger, *exec.Command, log, namespace)
	} else if exec.Script != nil {
		return script(ctx, logger, *exec.Script, log, namespace)
	}
	return nil
}

func command(ctx context.Context, logger logging.Logger, command v1alpha1.Command, log bool, namespace string) (_err error) {
	const operation = "CMD   "
	var output CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(color.BoldGreen.Sprint(operation), "DONE")
		} else {
			logger.Log(color.BoldRed.Sprint(operation), "ERROR", _err)
		}
	}()
	if log {
		defer func() {
			if out := output.Out(); out != "" {
				logger.Log(color.BoldFgCyan.Sprint("STDOUT"), "LOGS...\n"+out)
			}
			if err := output.Err(); err != "" {
				logger.Log(color.BoldFgCyan.Sprint("STDERR"), "LOGS...\n"+err)
			}
		}()
	} else {
		logger.Log(color.BoldFgCyan.Sprint("STDXXX"), "suppressed logs")
	}
	cmd := exec.CommandContext(ctx, command.EntryPoint, command.Args...) //nolint:gosec
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
	logger.Log(color.BoldFgCyan.Sprint(operation), cmd, "RUNNING...")
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr
	return cmd.Run()
}

func script(ctx context.Context, logger logging.Logger, script v1alpha1.Script, log bool, namespace string) (_err error) {
	const operation = "SCRIPT"
	var output CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(color.BoldGreen.Sprint(operation), "DONE")
		} else {
			logger.Log(color.BoldRed.Sprint(operation), "ERROR", _err)
		}
	}()
	if log {
		defer func() {
			if out := output.Out(); out != "" {
				logger.Log(color.BoldFgCyan.Sprint("STDOUT"), "LOGS...\n"+out)
			}
			if err := output.Err(); err != "" {
				logger.Log(color.BoldFgCyan.Sprint("STDERR"), "LOGS...\n"+err)
			}
		}()
	} else {
		logger.Log(color.BoldFgCyan.Sprint("STDXXX"), "suppressed logs")
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
	logger.Log(color.BoldFgCyan.Sprint(operation), "RUNNING...")
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr
	return cmd.Run()
}
