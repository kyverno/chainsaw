package script

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
)

type operation struct {
	script    v1alpha1.Script
	namespace string
}

func New(script v1alpha1.Script, namespace string) *operation {
	return &operation{
		script:    script,
		namespace: namespace,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	logger := logging.FromContext(ctx)
	var output internal.CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(logging.Script, color.BoldGreen, "DONE")
		} else {
			logger.Log(logging.Script, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	if !o.script.SkipLogOutput {
		defer func() {
			if out := output.Out(); out != "" {
				logger.Log(logging.Stdout, color.BoldFgCyan, "LOGS...\n"+out)
			}
			if err := output.Err(); err != "" {
				logger.Log(logging.Stdout, color.BoldFgCyan, "LOGS...\n"+err)
			}
		}()
	} else {
		logger.Log(logging.Std___, color.BoldYellow, "suppressed logs")
	}
	cmd := exec.CommandContext(ctx, "sh", "-c", o.script.Content) //nolint:gosec
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAMESPACE=%s", o.namespace))
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	logger.Log(logging.Script, color.BoldFgCyan, "RUNNING...")
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	cmdErr := cmd.Run()
	if o.script.Check.Value == nil {
		return cmdErr
	} else {
		actual := map[string]interface{}{
			"error":  nil,
			"stdout": output.Out(),
			"stderr": output.Err(),
		}
		if cmdErr != nil {
			actual["error"] = cmdErr.Error()
		}
		errs, err := assert.Validate(ctx, o.script.Check.Value, actual, nil)
		if err != nil {
			return err
		}
		return errs.ToAggregate()
	}
}
