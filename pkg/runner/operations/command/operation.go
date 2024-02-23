package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/env"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
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

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (_err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Command, _err)
	}()
	cmd, cancel, err := o.createCommand(ctx)
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		return err
	}
	internal.LogStart(logger, logging.Command, logging.Section("COMMAND", cmd.String()))
	return o.execute(ctx, bindings, cmd)
}

func (o *operation) createCommand(ctx context.Context) (*exec.Cmd, context.CancelFunc, error) {
	var cancel context.CancelFunc
	args := env.Expand(map[string]string{"NAMESPACE": o.namespace}, o.command.Args...)
	cmd := exec.CommandContext(ctx, o.command.Entrypoint, args...) //nolint:gosec
	env := os.Environ()
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory (%w)", err)
	}
	env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	env = append(env, fmt.Sprintf("NAMESPACE=%s", o.namespace))
	if o.cfg != nil {
		f, err := os.CreateTemp(o.basePath, "chainsaw-kubeconfig-")
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
		fmt.Println(f.Name())
		fmt.Println(o.basePath)
		env = append(env, fmt.Sprintf("KUBECONFIG=%s", filepath.Join(cwd, path)))
	}
	cmd.Env = env
	cmd.Dir = o.basePath
	return cmd, cancel, nil
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, cmd *exec.Cmd) error {
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
	bindings = bindings.Register("$stdout", binding.NewBinding(output.Out()))
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
