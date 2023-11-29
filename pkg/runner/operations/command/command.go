package command

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
	command   v1alpha1.Command
	basePath  string
	namespace string
}

func New(command v1alpha1.Command, basePath string, namespace string) operations.Operation {
	return &operation{
		command:   command,
		basePath:  basePath,
		namespace: namespace,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	logger := logging.FromContext(ctx)
	var output internal.CommandOutput
	defer func() {
		if _err == nil {
			logger.Log(logging.Command, logging.DoneStatus, color.BoldGreen)
		} else {
			logger.Log(logging.Command, logging.ErrorStatus, color.BoldRed, logging.ErrSection(_err))
		}
	}()
	if !o.command.SkipLogOutput {
		defer func() {
			if sections := output.Sections(); len(sections) != 0 {
				logger.Log(logging.Command, logging.LogStatus, color.BoldFgCyan, sections...)
			}
		}()
	}
	args := expand(map[string]string{"NAMESPACE": o.namespace}, o.command.Args...)
	cmd := exec.CommandContext(ctx, o.command.Entrypoint, args...) //nolint:gosec
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
	logger.Log(logging.Command, logging.RunStatus, color.BoldFgCyan, logging.Section("COMMAND", cmd.String()))
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	if err := cmd.Run(); o.command.Check == nil || o.command.Check.Value == nil {
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
		if errs, err := check.Check(ctx, nil, bindings, o.command.Check); err != nil {
			return err
		} else {
			return errs.ToAggregate()
		}
	}
}

func expand(env map[string]string, in ...string) []string {
	var args []string
	for _, arg := range in {
		expanded := os.Expand(arg, func(key string) string {
			expanded := env[key]
			if expanded == "" {
				expanded = os.Getenv(key)
			}
			return expanded
		})
		if expanded != "" {
			args = append(args, expanded)
		} else {
			args = append(args, arg)
		}
	}
	return args
}
