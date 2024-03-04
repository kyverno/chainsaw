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
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

type operation struct {
	client     client.Client
	base       unstructured.Unstructured
	namespacer namespacer.Namespacer
	cleaner    cleanup.Cleaner
	template   bool
	expect     []v1alpha1.Expectation
}

func New(
	client client.Client,
	obj unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	cleaner cleanup.Cleaner,
	template bool,
	expect []v1alpha1.Expectation,
) operations.Operation {
	return &operation{
		client:     client,
		base:       obj,
		namespacer: namespacer,
		cleaner:    cleaner,
		template:   template,
		expect:     expect,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (outputs operations.Outputs, err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Create, err)
	}()
	selfModifier := v1alpha1.Any{
		Value: obj.UnstructuredContent(),
	}
	if merged, err := mutate.Merge(ctx, obj, bindings, selfModifier); err != nil {
		return nil, err
	} else {
		obj = merged
	}
	if err := internal.ApplyNamespacer(o.namespacer, &obj); err != nil {
		return nil, err
	}
	internal.LogStart(logger, logging.Create)
	return nil, o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	var lastErr error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		lastErr = o.tryCreateResource(ctx, bindings, obj)
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

// TODO: could be replaced by checking the already exists error
func (o *operation) tryCreateResource(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.ObjectKey(&obj), &actual)
	if err == nil {
		return errors.New("the resource already exists in the cluster")
	}
	if kerrors.IsNotFound(err) {
		return o.createResource(ctx, bindings, obj)
	}
	return err
}

func (o *operation) createResource(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	err := o.client.Create(ctx, &obj)
	if err == nil && o.cleaner != nil {
		o.cleaner(obj, o.client)
	}
	return o.handleCheck(ctx, bindings, obj, err)
}

func (o *operation) handleCheck(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured, err error) error {
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if matched, err := check.Expectations(ctx, obj, bindings, o.expect...); matched {
		return err
	}
	return err
}
