package apply

import (
	"context"

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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
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
		internal.LogEnd(logger, logging.Apply, err)
	}()
	if o.template {
		template := v1alpha1.Any{
			Value: obj.UnstructuredContent(),
		}
		if merged, err := mutate.Merge(ctx, obj, bindings, template); err != nil {
			return nil, err
		} else {
			obj = merged
		}
	}
	if err := internal.ApplyNamespacer(o.namespacer, &obj); err != nil {
		return nil, err
	}
	internal.LogStart(logger, logging.Apply)
	return nil, o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	var lastErr error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		lastErr = o.tryApplyResource(ctx, bindings, obj)
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

func (o *operation) tryApplyResource(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.ObjectKey(&obj), &actual)
	if err == nil {
		return o.updateResource(ctx, bindings, &actual, obj)
	}
	if kerrors.IsNotFound(err) {
		return o.createResource(ctx, bindings, obj)
	}
	return err
}

func (o *operation) updateResource(ctx context.Context, bindings binding.Bindings, actual *unstructured.Unstructured, obj unstructured.Unstructured) error {
	patched, err := client.PatchObject(actual, &obj)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(patched)
	if err != nil {
		return err
	}
	return o.handleCheck(ctx, bindings, obj, o.client.Patch(ctx, actual, ctrlclient.RawPatch(types.MergePatchType, bytes)))
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
