package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	opapply "github.com/kyverno/chainsaw/pkg/engine/operations/apply"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type applyAction struct {
	op       v1alpha1.Apply
	resource unstructured.Unstructured
	cleaner  cleaner.CleanerCollector
}

func (o applyAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
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
		op := opapply.New(
			tc.Compilers(),
			client,
			o.resource,
			tc.Namespacer(),
			getCleanerOrNil(o.cleaner, tc),
			tc.Templating(),
			o.op.Expect,
			o.op.Outputs,
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Apply.Duration)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func applyOperation(ctx context.Context, tc enginecontext.TestContext, cleaner cleaner.CleanerCollector, op v1alpha1.Apply) ([]Operation, error) {
	resources, err := fileRefOrResource(ctx, op.ActionResourceRef, tc.BasePath(), tc.Compilers(), tc.Bindings())
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, applyAction{
			op:       op,
			resource: resource,
			cleaner:  cleaner,
		})
	}
	return ops, nil
}
