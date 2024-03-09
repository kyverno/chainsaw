package script

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/client-go/rest"
)

type operation struct {
	script    v1alpha1.Script
	basePath  string
	namespace string
	cfg       *rest.Config
}

func New(
	script v1alpha1.Script,
	basePath string,
	namespace string,
	cfg *rest.Config,
) operations.Operation {
	return &operation{
		script:    script,
		basePath:  basePath,
		namespace: namespace,
		cfg:       cfg,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (_ operations.Outputs, _err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Script, _err)
	}()
	cmd, cancel, _err := o.createCommand(ctx, bindings)
	if cancel != nil {
		defer cancel()
	}
	if _err != nil {
		return nil, _err
	}
	internal.LogStart(logger, logging.Script, logging.Section("COMMAND", cmd.String()))
	return o.execute(ctx, bindings, cmd)
}

func (o *operation) createCommand(ctx context.Context, bindings binding.Bindings) (*exec.Cmd, context.CancelFunc, error) {
	_, envs, err := internal.RegisterEnvs(ctx, o.namespace, bindings, o.script.Env...)
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
	cmd := exec.CommandContext(ctx, "sh", "-c", o.script.Content) //nolint:gosec
	cmd.Env = env
	cmd.Dir = o.basePath
	return cmd, cancel, nil
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, cmd *exec.Cmd) (_outputs operations.Outputs, _err error) {
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
	bindings = apibindings.RegisterNamedBinding(ctx, bindings, "stdout", output.Out())
	bindings = apibindings.RegisterNamedBinding(ctx, bindings, "stderr", output.Err())
	if err == nil {
		bindings = apibindings.RegisterNamedBinding(ctx, bindings, "error", nil)
	} else {
		bindings = apibindings.RegisterNamedBinding(ctx, bindings, "error", err.Error())
	}
	defer func(bindings binding.Bindings) {
		var outputs operations.Outputs
		if _err == nil {
			for _, output := range o.script.Outputs {
				if output.Match != nil && output.Match.Value != nil {
					if errs, err := check.Check(ctx, nil, nil, output.Match); err != nil {
						_err = err
						return
					} else if len(errs) != 0 {
						continue
					}
				}
				name, value, err := apibindings.ResolveBinding(ctx, bindings, nil, output.Binding)
				if err != nil {
					_err = err
					return
				}
				bindings = apibindings.RegisterNamedBinding(ctx, bindings, name, value)
				if outputs == nil {
					outputs = operations.Outputs{}
				}
				outputs[name] = value
			}
			_outputs = outputs
		}
	}(bindings)
	if o.script.Check == nil || o.script.Check.Value == nil {
		return nil, err
	}
	if errs, err := check.Check(ctx, nil, bindings, o.script.Check); err != nil {
		return nil, err
	} else {
		return nil, errs.ToAggregate()
	}
}
