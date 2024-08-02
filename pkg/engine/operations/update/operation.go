package update

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

const op = logging.Update

type operation struct {
	client     client.Client
	base       unstructured.Unstructured
	namespacer namespacer.Namespacer
	template   bool
	expect     []v1alpha1.Expectation
	outputs    []v1alpha1.Output
}

func New(
	client client.Client,
	obj unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	template bool,
	expect []v1alpha1.Expectation,
	outputs []v1alpha1.Output,
) operations.Operation {
	return &operation{
		client:     client,
		base:       obj,
		namespacer: namespacer,
		template:   template,
		expect:     expect,
		outputs:    outputs,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (_ outputs.Outputs, _err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, op, _err)
	}()
	if o.template {
		template := v1alpha1.Any{
			Value: obj.UnstructuredContent(),
		}
		if merged, err := templating.TemplateAndMerge(ctx, obj, bindings, template); err != nil {
			return nil, err
		} else {
			obj = merged
		}
	}
	if err := internal.ApplyNamespacer(o.namespacer, o.client, &obj); err != nil {
		return nil, err
	}
	internal.LogStart(logger, op)
	return o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	var lastErr error
	var outputs outputs.Outputs
	err := wait.PollUntilContextCancel(ctx, client.PollInterval, false, func(ctx context.Context) (bool, error) {
		outputs, lastErr = o.tryUpdateResource(ctx, bindings, obj)
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
func (o *operation) tryUpdateResource(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.Key(&obj), &actual)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, errors.New("the resource does not exist in the cluster")
		}
		return nil, err
	}
	obj.SetResourceVersion(actual.GetResourceVersion())
	return o.updateResource(ctx, bindings, obj)
}

func (o *operation) updateResource(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	err := o.client.Update(ctx, &obj)
	return o.handleCheck(ctx, bindings, obj, err)
}

func (o *operation) handleCheck(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured, err error) (_outputs outputs.Outputs, _err error) {
	if err == nil {
		bindings = apibindings.RegisterBinding(ctx, bindings, "error", nil)
	} else {
		bindings = apibindings.RegisterBinding(ctx, bindings, "error", err.Error())
	}
	defer func(bindings binding.Bindings) {
		if _err == nil {
			outputs, err := outputs.Process(ctx, bindings, obj.UnstructuredContent(), o.outputs...)
			if err != nil {
				_err = err
				return
			}
			_outputs = outputs
		}
	}(bindings)
	if matched, err := checks.Expect(ctx, obj, bindings, o.expect...); matched {
		return nil, err
	}
	return nil, err
}
