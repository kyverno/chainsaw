package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	opdelete "github.com/kyverno/chainsaw/pkg/engine/operations/delete"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type deleteAction struct {
	basePath   string
	namespacer namespacer.Namespacer
	op         v1alpha1.Delete
	resource   unstructured.Unstructured
}

func (o deleteAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	contextData := enginecontext.ContextData{
		BasePath:            o.basePath,
		Cluster:             o.op.Cluster,
		Clusters:            o.op.Clusters,
		DeletionPropagation: o.op.DeletionPropagationPolicy,
		Templating:          o.op.Template,
		Timeouts:            &v1alpha1.Timeouts{Delete: o.op.Timeout},
	}
	if tc, err := enginecontext.SetupContextAndBindings(tc, contextData, o.op.Bindings...); err != nil {
		return nil, err
	} else if _, client, err := tc.CurrentClusterClient(); err != nil {
		return nil, err
	} else {
		op := opdelete.New(
			tc.Compilers(),
			client,
			o.resource,
			o.namespacer,
			tc.Templating(),
			tc.DeletionPropagation(),
			o.op.Expect...,
		)
		ctx, cancel := context.WithTimeout(ctx, tc.Timeouts().Delete.Duration)
		defer cancel()
		return op.Exec(ctx, tc.Bindings())
	}
}

func deleteOperation(compilers compilers.Compilers, basePath string, namespacer namespacer.Namespacer, bindings apis.Bindings, op v1alpha1.Delete) ([]Operation, error) {
	ref := v1alpha1.ActionResourceRef{
		FileRef: v1alpha1.FileRef{
			File: op.File,
		},
	}
	if op.Ref != nil {
		var resource unstructured.Unstructured
		resource.SetAPIVersion(string(op.Ref.APIVersion))
		resource.SetKind(string(op.Ref.Kind))
		resource.SetName(string(op.Ref.Name))
		resource.SetNamespace(string(op.Ref.Namespace))
		resource.SetLabels(op.Ref.Labels)
		ref.Resource = &resource
	}
	resources, err := fileRefOrResource(context.TODO(), ref, basePath, compilers, bindings)
	if err != nil {
		return nil, err
	}
	var ops []Operation
	for i := range resources {
		resource := resources[i]
		ops = append(ops, deleteAction{
			basePath:   basePath,
			namespacer: namespacer,
			op:         op,
			resource:   resource,
		})
	}
	return ops, nil
}
