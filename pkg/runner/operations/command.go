package operations

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
)

func Command(ctx context.Context, logger logging.Logger, command v1alpha1.Command, namespace string) (*CommandOutput, error) {
	logger = logger.WithName("CMD   ")
	logger.Log(command.Command, "...")
	args := strings.Fields(command.Command)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint:gosec
	var output *CommandOutput
	if !command.SkipLogOutput {
		output = &CommandOutput{}
		cmd.Stdout = &output.stdout
		cmd.Stderr = &output.stderr
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAMESPACE=%s", namespace))
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	return output, cmd.Run()
}
