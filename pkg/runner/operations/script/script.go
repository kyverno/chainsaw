package script

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno/ext/output/color"
)

type operation struct {
	script    v1alpha1.Script
	basePath  string
	namespace string
}

func New(script v1alpha1.Script, basePath string, namespace string) operations.Operation {
	return &operation{
		script:    script,
		basePath:  basePath,
		namespace: namespace,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	logger := logging.FromContext(ctx)
	var output internal.CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(logging.Script, logging.DoneStatus, color.BoldGreen)
		} else {
			logger.Log(logging.Script, logging.ErrorStatus, color.BoldRed, logging.ErrSection(_err))
		}
	}()
	if !o.script.SkipLogOutput {
		defer func() {
			logger.Log(logging.Script, logging.LogStatus, color.BoldFgCyan, output.Sections()...)
		}()
	}
	cmd := exec.CommandContext(ctx, "sh", "-c", o.script.Content) //nolint:gosec
	env := os.Environ()
	if cwd, err := os.Getwd(); err != nil {
		return fmt.Errorf("failed to get current working directory (%w)", err)
	} else {
		env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	}
	env = append(env, fmt.Sprintf("NAMESPACE=%s", o.namespace))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	cmd.Dir = o.basePath
	logger.Log(logging.Script, logging.RunStatus, color.BoldFgCyan, logging.Section("COMMAND", cmd.String()))
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	if err := cmd.Run(); o.script.Check == nil || o.script.Check.Value == nil {
		return err
	} else {
		bindings := binding.NewBindings()
		if err == nil {
			bindings = bindings.Register("$error", binding.NewBinding(nil))
		} else {
			bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
		}
		bindings = bindings.Register("$stdout", binding.NewBinding(output.Out()))
		bindings = bindings.Register("$stderr", binding.NewBinding(output.Err()))
		if errs, err := check.Check(ctx, nil, bindings, o.script.Check); err != nil {
			return err
		} else {
			return errs.ToAggregate()
		}
	}
}
