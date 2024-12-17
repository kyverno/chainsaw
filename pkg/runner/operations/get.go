package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/kubectl"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	opcommand "github.com/kyverno/chainsaw/pkg/engine/operations/command"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

type getAction struct {
	basePath   string
	namespacer namespacer.Namespacer
	op         v1alpha1.Get
}

func (o getAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	ns := ""
	if o.namespacer != nil {
		ns = o.namespacer.GetNamespace()
	}
	contextData := enginecontext.ContextData{
		BasePath: o.basePath,
		Cluster:  o.op.Cluster,
		Clusters: o.op.Clusters,
		Timeouts: &v1alpha1.Timeouts{Exec: o.op.Timeout},
	}
	if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
		return nil, err
	} else if config, client, err := tc.CurrentClusterClient(); err != nil {
		return nil, err
	} else {
		entrypoint, args, err := kubectl.Get(ctx, tc.Compilers(), client, tc.Bindings(), &o.op)
		if err != nil {
			return nil, err
		}
		op := opcommand.New(
			tc.Compilers(),
			v1alpha1.Command{
				ActionClusters: o.op.ActionClusters,
				ActionTimeout:  o.op.ActionTimeout,
				Entrypoint:     entrypoint,
				Args:           args,
			},
			o.basePath,
			ns,
			config,
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Exec.Duration)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func getOperation(basePath string, namespacer namespacer.Namespacer, op v1alpha1.Get) Operation {
	return getAction{
		basePath:   basePath,
		namespacer: namespacer,
		op:         op,
	}
}
