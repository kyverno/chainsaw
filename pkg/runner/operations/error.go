package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	operror "github.com/kyverno/chainsaw/pkg/engine/operations/error"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type errorAction struct {
	op       v1alpha1.Error
	resource unstructured.Unstructured
}

func (o errorAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	contextData := enginecontext.ContextData{
		Cluster:    o.op.Cluster,
		Clusters:   o.op.Clusters,
		Templating: o.op.Template,
		Timeouts:   &v1alpha1.Timeouts{Error: o.op.Timeout},
	}
	if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, o.op.Bindings...); err != nil {
		return nil, err
	} else if _, client, err := tc.CurrentClusterClient(); err != nil {
		return nil, err
	} else {
		op := operror.New(
			tc.Compilers(),
			client,
			o.resource,
			tc.Namespacer(),
			tc.Templating(),
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Error)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func errorOperation(ctx context.Context, tc enginecontext.TestContext, op v1alpha1.Error) ([]Operation, error) {
	resources, err := fileRefOrCheck(ctx, op.ActionCheckRef, tc.BasePath(), tc.Compilers(), tc.Bindings())
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, errorAction{
			op:       op,
			resource: resource,
		})
	}
	return ops, nil
}
