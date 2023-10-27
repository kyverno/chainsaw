package operations

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
)

func Exec(ctx context.Context, logger logging.Logger, exec v1alpha1.Exec, namespace string) (CommandOutput, error) {
	if exec.Command != nil {
		return command(ctx, logger, *exec.Command, namespace)
	} else if exec.Script != nil {
		return script(ctx, logger, *exec.Script, namespace)
	}
	return CommandOutput{}, nil
}

func command(ctx context.Context, logger logging.Logger, command v1alpha1.Command, namespace string) (CommandOutput, error) {
	cmd := exec.CommandContext(ctx, command.EntryPoint, command.Args...) //nolint:gosec
	logger = logger.WithName("CMD   ")
	logger.Log(cmd, "...")
	var output CommandOutput
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr
	cwd, err := os.Getwd()
	if err != nil {
		return output, fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAMESPACE=%s", namespace))
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	return output, cmd.Run()
}

func script(ctx context.Context, logger logging.Logger, script v1alpha1.Script, namespace string) (CommandOutput, error) {
	logger = logger.WithName("SCRIPT")
	logger.Log("executing...")
	cmd := exec.CommandContext(ctx, "sh", "-c", script.Content) //nolint:gosec
	var output CommandOutput
	cmd.Stdout = &output.stdout
	cmd.Stderr = &output.stderr
	cwd, err := os.Getwd()
	if err != nil {
		return output, fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAMESPACE=%s", namespace))
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	return output, cmd.Run()
}
