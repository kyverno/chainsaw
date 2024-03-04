package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	mutation "github.com/kyverno/chainsaw/pkg/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	environment "github.com/kyverno/chainsaw/pkg/runner/env"
	"github.com/kyverno/chainsaw/pkg/runner/functions"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/client-go/rest"
)

type operation struct {
	command   v1alpha1.Command
	basePath  string
	namespace string
	cfg       *rest.Config
}

func New(
	command v1alpha1.Command,
	basePath string,
	namespace string,
	cfg *rest.Config,
) operations.Operation {
	return &operation{
		command:   command,
		basePath:  basePath,
		namespace: namespace,
		cfg:       cfg,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (outputs operations.Outputs, _err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Command, _err)
	}()
	cmd, cancel, err := o.createCommand(ctx, bindings)
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		return nil, err
	}
	internal.LogStart(logger, logging.Command, logging.Section("COMMAND", cmd.String()))
	return o.execute(ctx, bindings, cmd)
}

func (o *operation) createCommand(ctx context.Context, bindings binding.Bindings) (*exec.Cmd, context.CancelFunc, error) {
	maps, envs, err := internal.RegisterEnvs(ctx, o.namespace, bindings, o.command.Env...)
	if err != nil {
		return nil, nil, err
	}
	env := os.Environ()
	env = append(env, envs...)
	var cancel context.CancelFunc
	if o.cfg != nil {
		f, err := os.CreateTemp("", "chainsaw-kubeconfig-")
		if err != nil {
			return nil, nil, err
		}
		path := f.Name()
		cancel = func() {
			err := os.Remove(path)
			if err != nil {
				logger := internal.GetLogger(ctx, nil)
				logger.Log(logging.Script, logging.ErrorStatus, color.BoldYellow, logging.ErrSection(err))
			}
		}
		defer f.Close()
		if err := restutils.Save(o.cfg, f); err != nil {
			return nil, cancel, err
		}
		env = append(env, fmt.Sprintf("KUBECONFIG=%s", path))
	}
	args := environment.Expand(maps, o.command.Args...)
	cmd := exec.CommandContext(ctx, o.command.Entrypoint, args...) //nolint:gosec
	cmd.Env = env
	cmd.Dir = o.basePath
	return cmd, cancel, nil
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, cmd *exec.Cmd) (_outputs operations.Outputs, _err error) {
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
	bindings = bindings.Register("$stdout", binding.NewBinding(output.Out()))
	bindings = bindings.Register("$stderr", binding.NewBinding(output.Err()))
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	defer func() {
		var outputs operations.Outputs
		if _err == nil {
			for _, output := range o.command.Outputs {
				if err := output.CheckName(); err != nil {
					_err = err
					return
				}
				patched, err := mutation.Mutate(ctx, nil, mutation.Parse(ctx, output.Value.Value), nil, bindings, template.WithFunctionCaller(functions.Caller))
				if err != nil {
					_err = err
					return
				}
				if outputs == nil {
					outputs = operations.Outputs{}
				}
				outputs[output.Name] = binding.NewBinding(patched)
				bindings = bindings.Register("$"+output.Name, outputs[output.Name])
			}
			_outputs = outputs
		}
	}()
	if o.command.Check == nil || o.command.Check.Value == nil {
		return nil, err
	}
	if errs, err := check.Check(ctx, nil, bindings, o.command.Check); err != nil {
		return nil, err
	} else {
		return nil, errs.ToAggregate()
	}
}
