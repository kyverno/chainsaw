package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	apibindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/logging"
	environment "github.com/kyverno/chainsaw/pkg/utils/env"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/client-go/rest"
)

type operation struct {
	compilers compilers.Compilers
	command   v1alpha1.Command
	basePath  string
	namespace string
	cfg       *rest.Config
}

func New(
	compilers compilers.Compilers,
	command v1alpha1.Command,
	basePath string,
	namespace string,
	cfg *rest.Config,
) operations.Operation {
	return &operation{
		compilers: compilers,
		command:   command,
		basePath:  basePath,
		namespace: namespace,
		cfg:       cfg,
	}
}

func (o *operation) Exec(ctx context.Context, bindings apis.Bindings) (_ outputs.Outputs, _err error) {
	if bindings == nil {
		bindings = apis.NewBindings()
	}
	defer func() {
		internal.LogEnd(ctx, logging.Command, nil, _err)
	}()
	cmd, cancel, err := o.createCommand(ctx, bindings)
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		return nil, err
	}
	var logOpts []fmt.Stringer
	if !o.command.SkipCommandOutput {
		logOpts = append(logOpts, logging.Section("COMMAND", cmd.String()))
	}
	internal.LogStart(ctx, logging.Command, nil, logOpts...)
	return o.execute(ctx, bindings, cmd)
}

func (o *operation) createCommand(ctx context.Context, bindings apis.Bindings) (*exec.Cmd, context.CancelFunc, error) {
	maps, envs, err := internal.RegisterEnvs(ctx, o.compilers, o.namespace, bindings, o.command.Env...)
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
				logging.Log(ctx, logging.Script, logging.WarnStatus, nil, color.BoldYellow, logging.ErrSection(err))
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
	basePath := o.basePath
	if o.command.WorkDir != nil {
		if filepath.IsAbs(*o.command.WorkDir) {
			basePath = *o.command.WorkDir
		} else {
			basePath = filepath.Join(basePath, *o.command.WorkDir)
		}
	}
	cmd.Dir = basePath
	return cmd, cancel, nil
}

func (o *operation) execute(ctx context.Context, bindings apis.Bindings, cmd *exec.Cmd) (_outputs outputs.Outputs, _err error) {
	var output internal.CommandOutput
	if !o.command.SkipLogOutput {
		defer func() {
			if sections := output.Sections(); len(sections) != 0 {
				logging.Log(ctx, logging.Command, logging.LogStatus, nil, color.BoldFgCyan, sections...)
			}
		}()
	}
	cmd.Stdout = &output.Stdout
	cmd.Stderr = &output.Stderr
	err := cmd.Run()
	bindings = apibindings.RegisterBinding(bindings, "stdout", output.Out())
	bindings = apibindings.RegisterBinding(bindings, "stderr", output.Err())
	if err == nil {
		bindings = apibindings.RegisterBinding(bindings, "error", nil)
	} else {
		bindings = apibindings.RegisterBinding(bindings, "error", err.Error())
	}
	defer func(bindings apis.Bindings) {
		if _err == nil {
			outputs, err := outputs.Process(ctx, o.compilers, bindings, nil, o.command.Outputs...)
			if err != nil {
				_err = err
				return
			}
			_outputs = outputs
		}
	}(bindings)
	if o.command.Check == nil || o.command.Check.IsNil() {
		return nil, err
	}
	if errs, err := checks.Check(ctx, o.compilers, nil, bindings, o.command.Check); err != nil {
		return nil, err
	} else {
		return nil, errs.ToAggregate()
	}
}
