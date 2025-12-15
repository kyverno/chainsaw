package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	opscript "github.com/kyverno/chainsaw/pkg/engine/operations/script"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

type scriptAction struct {
	op v1alpha1.Script
}

func (o scriptAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	ns := ""
	if namespacer := tc.Namespacer(); namespacer != nil {
		ns = namespacer.GetNamespace()
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
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Exec)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func scriptOperation(op v1alpha1.Script) Operation {
	return scriptAction{
		op: op,
	}
}
