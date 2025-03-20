package apply

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	bindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/wait"
)

type operation struct {
	compilers  compilers.Compilers
	client     client.Client
	base       unstructured.Unstructured
	namespacer namespacer.Namespacer
	cleaner    cleaner.CleanerCollector
	template   bool
	expect     []v1alpha1.Expectation
	outputs    []v1alpha1.Output
}

func New(
	compilers compilers.Compilers,
	client client.Client,
	obj unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	cleaner cleaner.CleanerCollector,
	template bool,
	expect []v1alpha1.Expectation,
	outputs []v1alpha1.Output,
) operations.Operation {
	return &operation{
		compilers:  compilers,
		client:     client,
		base:       obj,
		namespacer: namespacer,
		cleaner:    cleaner,
		template:   template,
		expect:     expect,
		outputs:    outputs,
	}
}

func (o *operation) Exec(ctx context.Context, tc apis.Bindings) (_ outputs.Outputs, _err error) {
	if tc == nil {
		tc = apis.NewBindings()
	}
	obj := o.base
	defer func() {
		internal.LogEnd(ctx, logging.Apply, &obj, _err)
	}()
	if o.template {
		template := v1alpha1.NewProjection(obj.UnstructuredContent())
		if merged, err := templating.TemplateAndMerge(ctx, o.compilers, obj, tc, template); err != nil {
			return nil, err
		} else {
			obj = merged
		}
	}
	if err := internal.ApplyNamespacer(o.namespacer, o.client, &obj); err != nil {
		return nil, err
	}
	internal.LogStart(ctx, logging.Apply, &obj)
	return o.execute(ctx, tc, obj)
}

func (o *operation) execute(ctx context.Context, tc apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	var lastErr error
	var outputs outputs.Outputs
	err := wait.PollUntilContextCancel(ctx, client.PollInterval, false, func(ctx context.Context) (bool, error) {
		outputs, lastErr = o.tryApplyResource(ctx, tc, obj)
		// TODO: determine if the error can be retried
		return lastErr == nil, nil
	})
	if err == nil {
		return outputs, nil
	}
	if lastErr != nil {
		return outputs, lastErr
	}
	return outputs, err
}

func (o *operation) tryApplyResource(ctx context.Context, tc apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.Key(&obj), &actual)
	operations.ResetWarnings(ctx)
	if err == nil {
		return o.updateResource(ctx, tc, &actual, obj)
	}
	if kerrors.IsNotFound(err) {
		return o.createResource(ctx, tc, obj)
	}
	return nil, err
}

func (o *operation) updateResource(ctx context.Context, tc apis.Bindings, actual *unstructured.Unstructured, obj unstructured.Unstructured) (outputs.Outputs, error) {
	patched, err := client.PatchObject(actual, &obj)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(patched)
	if err != nil {
		return nil, err
	}
	return o.handleCheck(ctx, tc, obj, o.client.Patch(ctx, actual, client.RawPatch(types.MergePatchType, bytes)))
}

func (o *operation) createResource(ctx context.Context, tc apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	err := o.client.Create(ctx, &obj)
	if err == nil && o.cleaner != nil {
		o.cleaner.Add(o.client, &obj)
	}
	return o.handleCheck(ctx, tc, obj, err)
}

func (o *operation) handleCheck(ctx context.Context, tc apis.Bindings, obj unstructured.Unstructured, err error) (_outputs outputs.Outputs, _err error) {
	if err == nil {
		tc = bindings.RegisterBinding(tc, "error", nil)
	} else {
		tc = bindings.RegisterBinding(tc, "error", err.Error())
	}
	tc = operations.RegisterWarningsInBindings(ctx, tc)
	defer func(tc apis.Bindings) {
		if _err == nil {
			outputs, err := outputs.Process(ctx, o.compilers, tc, obj.UnstructuredContent(), o.outputs...)
			if err != nil {
				_err = err
				return
			}
			_outputs = outputs
		}
	}(tc)
	if matched, err := checks.Expect(ctx, o.compilers, obj, tc, o.expect...); matched {
		return nil, err
	}
	return nil, err
}
