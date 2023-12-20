package create

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

type operation struct {
	client     client.Client
	obj        unstructured.Unstructured
	namespacer namespacer.Namespacer
	cleaner    cleanup.Cleaner
	expect     []v1alpha1.Expectation
}

func New(client client.Client, obj unstructured.Unstructured, namespacer namespacer.Namespacer, cleaner cleanup.Cleaner, expect ...v1alpha1.Expectation) operations.Operation {
	return &operation{
		client:     client,
		obj:        obj,
		namespacer: namespacer,
		cleaner:    cleaner,
		expect:     expect,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := internal.GetLogger(ctx, &o.obj)
	defer func() {
		internal.LogEnd(logger, logging.Create, err)
	}()
	if err := internal.ApplyNamespacer(o.namespacer, &o.obj); err != nil {
		return err
	}
	internal.LogStart(logger, logging.Create)
	return o.execute(ctx, logger)
}

func (o *operation) execute(ctx context.Context, logger logging.Logger) error {
	var lastErr error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		lastErr = o.tryCreateResource(ctx)
		// TODO: determine if the error can be retried
		return lastErr == nil, nil
	})
	if err == nil {
		return nil
	}
	if lastErr != nil {
		return lastErr
	}
	return err
}

func (o *operation) tryCreateResource(ctx context.Context) error {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(o.obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.ObjectKey(&o.obj), &actual)
	if err == nil {
		return errors.New("the resource already exists in the cluster")
	}
	if kerrors.IsNotFound(err) {
		return o.createResource(ctx)
	}
	return err
}

func (o *operation) createResource(ctx context.Context) error {
	err := o.client.Create(ctx, &o.obj)
	if err == nil && o.cleaner != nil {
		o.cleaner(o.obj, o.client)
	}
	return o.handleCheck(ctx, err)
}

func (o *operation) handleCheck(ctx context.Context, err error) error {
	bindings := binding.NewBindings()
	bindings = bindings.Register("$client", binding.NewBinding(o.client))
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if matched, err := check.Expectations(ctx, o.obj, bindings, o.expect...); matched {
		return err
	}
	return err
}
