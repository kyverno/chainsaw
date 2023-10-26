package operations

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func ExecuteCommand(cmdSpec v1alpha1.Command, namespace string) error {
	var cmdStr string
	if cmdSpec.Script != "" {
		cmdStr = cmdSpec.Script
	} else {
		cmdStr = cmdSpec.Command
	}

	// If namespaced, add `--namespace` flag.
	if cmdSpec.Namespaced {
		cmdStr = fmt.Sprintf("%s --namespace %s", cmdStr, namespace)
	}

	// Handle timeout if provided
	var ctx context.Context = context.Background()
	var cancel context.CancelFunc

	if cmdSpec.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(cmdSpec.Timeout)*time.Second)
		defer cancel()
	}

	command := exec.CommandContext(ctx, "sh", "-c", cmdStr)

	// Handle output (if not skipping)
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
		fmt.Printf("Command output: %s\n", out) // i have to consider this as a log output or not
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
