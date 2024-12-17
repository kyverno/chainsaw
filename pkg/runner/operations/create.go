package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	opcreate "github.com/kyverno/chainsaw/pkg/engine/operations/create"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type createAction struct {
	basePath   string
	namespacer namespacer.Namespacer
	op         v1alpha1.Create
	resource   unstructured.Unstructured
	cleaner    cleaner.CleanerCollector
}

func (o createAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	contextData := enginecontext.ContextData{
		BasePath:   o.basePath,
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
		op := opcreate.New(
			tc.Compilers(),
			client,
			o.resource,
			o.namespacer,
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

func createOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, cleaner cleaner.CleanerCollector, bindings apis.Bindings, op v1alpha1.Create) ([]Operation, error) {
	resources, err := fileRefOrResource(context.TODO(), op.ActionResourceRef, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, createAction{
			basePath:   basePath,
			namespacer: namespacer,
			op:         op,
			resource:   resource,
			cleaner:    cleaner,
		})
	}
	return ops, nil
}
