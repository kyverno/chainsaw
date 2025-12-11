package create

import (
	"context"
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
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

func (o *operation) Exec(ctx context.Context, bindings apis.Bindings) (_ outputs.Outputs, _err error) {
	if bindings == nil {
		bindings = apis.NewBindings()
	}
	obj := o.base
	defer func() {
		internal.LogEnd(ctx, logging.Create, &obj, _err)
	}()
	if o.client == nil {
		return nil, errors.New("cluster client not set")
	}
	if o.template {
		template := v1alpha1.NewProjection(obj.UnstructuredContent())
		if merged, err := templating.TemplateAndMerge(ctx, o.compilers, obj, bindings, template); err != nil {
			return nil, err
		} else {
			obj = merged
		}
	}
	if err := internal.ApplyNamespacer(o.namespacer, o.client, &obj); err != nil {
		return nil, err
	}
	internal.LogStart(ctx, logging.Create, &obj)
	return o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	var lastErr error
	var outputs outputs.Outputs
	err := wait.PollUntilContextCancel(ctx, client.PollInterval, false, func(ctx context.Context) (bool, error) {
		outputs, lastErr = o.tryCreateResource(ctx, bindings, obj)
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

// TODO: could be replaced by checking the already exists error
func (o *operation) tryCreateResource(ctx context.Context, bindings apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.Key(&obj), &actual)
	if err == nil {
		return nil, errors.New("the resource already exists in the cluster")
	}
	if kerrors.IsNotFound(err) {
		return o.createResource(ctx, bindings, obj)
	}
	return nil, err
}

func (o *operation) createResource(ctx context.Context, bindings apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	err := o.client.Create(ctx, &obj)
	if err == nil && o.cleaner != nil {
		o.cleaner.Add(o.client, &obj)
	}
	return o.handleCheck(ctx, bindings, obj, err)
}

func (o *operation) handleCheck(ctx context.Context, bindings apis.Bindings, obj unstructured.Unstructured, err error) (_outputs outputs.Outputs, _err error) {
	if err == nil {
		bindings = apibindings.RegisterBinding(bindings, "error", nil)
	} else {
		bindings = apibindings.RegisterBinding(bindings, "error", err.Error())
	}
	defer func(bindings apis.Bindings) {
		if _err == nil {
			outputs, err := outputs.Process(ctx, o.compilers, bindings, obj.UnstructuredContent(), o.outputs...)
			if err != nil {
				_err = err
				return
			}
			_outputs = outputs
		}
	}(bindings)
	if matched, err := checks.Expect(ctx, o.compilers, obj, bindings, o.expect...); matched {
		return nil, err
	}
	return nil, err
}
