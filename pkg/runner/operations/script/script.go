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

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Script, err)
	}()
	cmd, err := o.createCommand(ctx)
	if err != nil {
		return err
	}
	internal.LogStart(logger, logging.Script, logging.Section("COMMAND", cmd.String()))
	return o.execute(ctx, cmd)
}

func (o *operation) createCommand(ctx context.Context) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", o.script.Content) //nolint:gosec
	env := os.Environ()
	if cwd, err := os.Getwd(); err != nil {
		return nil, fmt.Errorf("failed to get current working directory (%w)", err)
	} else {
		env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
		// TODO: won't work with multicluster support
		env = append(env, fmt.Sprintf("KUBECONFIG=%s/kubeconfig", cwd))
	}
	env = append(env, fmt.Sprintf("NAMESPACE=%s", o.namespace))
	cmd.Env = env
	cmd.Dir = o.basePath
	return cmd, nil
}

func (o *operation) execute(ctx context.Context, cmd *exec.Cmd) error {
	logger := internal.GetLogger(ctx, nil)
	var output internal.CommandOutput
	if !o.script.SkipLogOutput {
		defer func() {
			if sections := output.Sections(); len(sections) != 0 {
				logger.Log(logging.Script, logging.LogStatus, color.BoldFgCyan, sections...)
			}
		}()
	}
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	err := cmd.Run()
	if o.script.Check == nil || o.script.Check.Value == nil {
		return err
	}
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
