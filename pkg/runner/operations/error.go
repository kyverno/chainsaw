package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	operror "github.com/kyverno/chainsaw/pkg/engine/operations/error"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type errorAction struct {
	basePath   string
	namespacer namespacer.Namespacer
	op         v1alpha1.Error
	resource   unstructured.Unstructured
}

func (o errorAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	contextData := enginecontext.ContextData{
		BasePath:   o.basePath,
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
			o.namespacer,
			tc.Templating(),
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Error.Duration)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func errorOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Error) ([]Operation, error) {
	resources, err := fileRefOrCheck(context.TODO(), op.ActionCheckRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, errorAction{
			basePath:   basePath,
			namespacer: namespacer,
			op:         op,
			resource:   resource,
		})
	}
	return ops, nil
}
