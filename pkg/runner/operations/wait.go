package operations

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/kubectl"
	opcommand "github.com/kyverno/chainsaw/pkg/engine/operations/command"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type waitAction struct {
	op v1alpha1.Wait
}

func (o waitAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	ns := ""
	if namespacer := tc.Namespacer(); namespacer != nil {
		ns = namespacer.GetNamespace()
	}
	contextData := enginecontext.ContextData{
		Cluster:  o.op.Cluster,
		Clusters: o.op.Clusters,
		Timeouts: &v1alpha1.Timeouts{Exec: o.op.Timeout},
	}
	if tc, err := enginecontext.SetupContextAndBindings(tc, contextData); err != nil {
		return nil, err
	} else if config, client, err := tc.CurrentClusterClient(); err != nil {
		return nil, err
	} else {
		// make sure timeout is set to populate the command flag
		timeout := tc.Timeouts().Exec.Duration
		o.op.Timeout = &metav1.Duration{Duration: timeout}
		entrypoint, args, err := kubectl.Wait(ctx, tc.Compilers(), client, tc.Bindings(), &o.op)
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
			tc.BasePath(),
			ns,
			config,
		)
		// shift operation timeout
		ctx, cancel := context.WithTimeout(ctx, timeout+30*time.Second)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func waitOperation(op v1alpha1.Wait) Operation {
	return waitAction{
		op: op,
	}
}
