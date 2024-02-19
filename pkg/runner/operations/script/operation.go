package script

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
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
	bindings  binding.Bindings
	cfg       *rest.Config
}

func New(
	script v1alpha1.Script,
	basePath string,
	namespace string,
	bindings binding.Bindings,
	cfg *rest.Config,
) operations.Operation {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return &operation{
		script:    script,
		basePath:  basePath,
		namespace: namespace,
		bindings:  bindings,
		cfg:       cfg,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Script, err)
	}()
	cmd, cancel, err := o.createCommand(ctx)
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		return err
	}
	internal.LogStart(logger, logging.Script, logging.Section("COMMAND", cmd.String()))
	return o.execute(ctx, cmd)
}

func (o *operation) createCommand(ctx context.Context) (*exec.Cmd, context.CancelFunc, error) {
	var cancel context.CancelFunc
	cmd := exec.CommandContext(ctx, "sh", "-c", o.script.Content) //nolint:gosec
	env := os.Environ()
	if cwd, err := os.Getwd(); err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory (%w)", err)
	} else {
		env = append(env, fmt.Sprintf("PATH=%s/bin/:%s", cwd, os.Getenv("PATH")))
	}
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
		env = append(env, fmt.Sprintf("KUBECONFIG=%s", filepath.Base(path)))
	}
	cmd.Env = env
	cmd.Dir = o.basePath
	return cmd, cancel, nil
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
	bindings := o.bindings.Register("$stdout", binding.NewBinding(output.Out()))
	bindings = bindings.Register("$stderr", binding.NewBinding(output.Err()))
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if errs, err := check.Check(ctx, nil, bindings, o.script.Check); err != nil {
		return err
	} else {
		return errs.ToAggregate()
	}
}
