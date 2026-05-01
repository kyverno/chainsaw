package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	oppatch "github.com/kyverno/chainsaw/pkg/engine/operations/patch"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type patchAction struct {
	op       v1alpha1.Patch
	resource unstructured.Unstructured
}

func (o patchAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	contextData := enginecontext.ContextData{
		Cluster:    o.op.Cluster,
		Clusters:   o.op.Clusters,
		DryRun:     o.op.DryRun,
		Templating: o.op.Template,
		Timeouts:   &v1alpha1.Timeouts{Apply: o.op.Timeout},
	}
	if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, o.op.Bindings...); err != nil {
		return nil, err
	} else if err := prepareResource(o.resource, tc); err != nil {
		return nil, err
	} else if _, client, err := tc.CurrentClusterClient(); err != nil {
		return nil, err
	} else {
		op := oppatch.New(
			tc.Compilers(),
			client,
			o.resource,
			tc.Namespacer(),
			tc.Templating(),
			o.op.Expect,
			o.op.Outputs,
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Apply)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func patchOperation(ctx context.Context, tc enginecontext.TestContext, op v1alpha1.Patch) ([]Operation, error) {
	resources, err := fileRefOrResource(ctx, op.ActionResourceRef, tc.BasePath(), tc.Compilers(), tc.Bindings())
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, patchAction{
			op:       op,
			resource: resource,
		})
	}
	return ops, nil
}
