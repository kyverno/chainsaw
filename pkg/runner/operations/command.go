package operations

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func ExecuteCommand(ctx context.Context, cmdSpec v1alpha1.Command, namespace string) error {
	isScript := cmdSpec.Script != ""
	var cmdStr string
	if isScript {
		cmdStr = cmdSpec.Script
	} else {
		cmdStr = cmdSpec.Command
	}

	// Append namespace flag if necessary
	if cmdSpec.Namespaced && !isScript {
		cmdStr = fmt.Sprintf("%s --namespace %s", cmdStr, namespace)
	}

	// Create the appropriate command based on the input
	var command *exec.Cmd
	if isScript {
		// #nosec
		command = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	} else {
		args := strings.Fields(cmdStr)
		// #nosec
		command = exec.CommandContext(ctx, args[0], args[1:]...)
	}

	// Execute and handle output
	var out []byte
	var err error

	if !cmdSpec.SkipLogOutput {
		out, err = command.CombinedOutput()
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("command timed out after %d seconds", cmdSpec.Timeout)
			}
			return fmt.Errorf("command failed: %s\n%s", err, out)
		}
		// fmt.Printf("Command output: %s\n", out)
	} else {
		err = command.Run()
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("command timed out after %d seconds", cmdSpec.Timeout)
			}
			return err
		}
	}

	return nil
}
