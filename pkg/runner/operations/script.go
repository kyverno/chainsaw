package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	opscript "github.com/kyverno/chainsaw/pkg/engine/operations/script"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

type scriptAction struct {
	namespacer namespacer.Namespacer
	op         v1alpha1.Script
}

func (o scriptAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	ns := ""
	if o.namespacer != nil {
		ns = o.namespacer.GetNamespace()
	}
	contextData := enginecontext.ContextData{
		Cluster:  o.op.Cluster,
		Clusters: o.op.Clusters,
		Timeouts: &v1alpha1.Timeouts{Exec: o.op.Timeout},
	}
	if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, o.op.Bindings...); err != nil {
		return nil, err
	} else if config, _, err := tc.CurrentClusterClient(); err != nil {
		return nil, err
	} else {
		op := opscript.New(
			tc.Compilers(),
			o.op,
			tc.BasePath(),
			ns,
			config,
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Exec.Duration)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func scriptOperation(namespacer namespacer.Namespacer, op v1alpha1.Script) Operation {
	return scriptAction{
		namespacer: namespacer,
		op:         op,
	}
}
