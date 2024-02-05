package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/env"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno/ext/output/color"
)

type operation struct {
	command   v1alpha1.Command
	basePath  string
	namespace string
	bindings  binding.Bindings
}

func New(command v1alpha1.Command, basePath string, namespace string, bindings binding.Bindings) operations.Operation {
	return &operation{
		command:   command,
		basePath:  basePath,
		namespace: namespace,
		bindings:  bindings,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Command, _err)
	}()
	cmd, err := o.createCommand(ctx)
	if err != nil {
		return err
	}
	internal.LogStart(logger, logging.Command, logging.Section("COMMAND", cmd.String()))
	return o.execute(ctx, cmd)
}

func (o *operation) createCommand(ctx context.Context) (*exec.Cmd, error) {
	args := env.Expand(map[string]string{"NAMESPACE": o.namespace}, o.command.Args...)
	cmd := exec.CommandContext(ctx, o.command.Entrypoint, args...) //nolint:gosec
	env := os.Environ()
	if cwd, err := os.Getwd(); err != nil {
		return nil, fmt.Errorf("failed to get current working directory (%w)", err)
	} else {
		env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	}
	env = append(env, fmt.Sprintf("NAMESPACE=%s", o.namespace))
	// TODO
	// env = append(env, fmt.Sprintf("KUBECONFIG=%s/bin/:%s", cwd, os.Getenv("PATH")))
	cmd.Env = env
	cmd.Dir = o.basePath
	return cmd, nil
}

func (o *operation) execute(ctx context.Context, cmd *exec.Cmd) error {
	logger := logging.FromContext(ctx)
	var output internal.CommandOutput
	if !o.command.SkipLogOutput {
		defer func() {
			if sections := output.Sections(); len(sections) != 0 {
				logger.Log(logging.Command, logging.LogStatus, color.BoldFgCyan, sections...)
			}
		}()
	}
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	err := cmd.Run()
	if o.command.Check == nil || o.command.Check.Value == nil {
		return err
	}
	bindings := o.bindings.Register("$stdout", binding.NewBinding(output.Out()))
	bindings = bindings.Register("$stderr", binding.NewBinding(output.Err()))
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if errs, err := check.Check(ctx, nil, bindings, o.command.Check); err != nil {
		return err
	} else {
		return errs.ToAggregate()
	}
}
