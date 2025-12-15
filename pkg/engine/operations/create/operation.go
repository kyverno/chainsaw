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
		// Check if the error is retryable
		if lastErr != nil {
			// Conflict errors should be retried
			if kerrors.IsConflict(lastErr) {
				return false, nil
			}
			// Server timeout errors should be retried
			if kerrors.IsServerTimeout(lastErr) {
				return false, nil
			}
			// Too many requests errors should be retried
			if kerrors.IsTooManyRequests(lastErr) {
				return false, nil
			}
			// Service unavailable errors should be retried
			if kerrors.IsServiceUnavailable(lastErr) {
				return false, nil
			}
			// AlreadyExists error should not be retried as it's a permanent condition
			if kerrors.IsAlreadyExists(lastErr) {
				return false, lastErr
			}
			// Non-retryable error
			return false, lastErr
		}
		return true, nil
	})
	if err == nil {
		return outputs, nil
	}
	if lastErr != nil {
		return outputs, lastErr
	}
	return outputs, err
}

func (o *operation) tryCreateResource(ctx context.Context, bindings apis.Bindings, obj unstructured.Unstructured) (outputs.Outputs, error) {
	// First check if the resource exists
	key := client.Key(&obj)
	var existing unstructured.Unstructured
	existing.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())

	err := o.client.Get(ctx, key, &existing)
	if err != nil {
		// If there was an error other than NotFound, propagate it
		if !kerrors.IsNotFound(err) {
			return nil, err
		}

		// Resource doesn't exist, try to create it
		createErr := o.client.Create(ctx, &obj)
		if createErr == nil {
			// Resource created successfully
			if o.cleaner != nil {
				o.cleaner.Add(o.client, &obj)
			}
			return o.handleCheck(ctx, bindings, obj, nil)
		}

		// Check if the error matches any expectations
		if len(o.expect) > 0 {
			return o.handleCheck(ctx, bindings, obj, createErr)
		}

		// If the error is not AlreadyExists, propagate it
		if !kerrors.IsAlreadyExists(createErr) {
			return nil, createErr
		}

		// Resource already exists (race condition)
		return nil, errors.New("the resource already exists in the cluster")
	}

	// Resource already exists
	return nil, errors.New("the resource already exists in the cluster")
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
