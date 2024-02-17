package patch

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
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
	bindings   binding.Bindings
	template   bool
	expect     []v1alpha1.Expectation
}

func New(
	client client.Client,
	obj unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	bindings binding.Bindings,
	template bool,
	expect []v1alpha1.Expectation,
) operations.Operation {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return &operation{
		client:     client,
		base:       obj,
		namespacer: namespacer,
		bindings:   bindings,
		template:   template,
		expect:     expect,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Patch, err)
	}()
	selfModifier := v1alpha1.Any{
		Value: obj.UnstructuredContent(),
	}
	if merged, err := mutate.Merge(ctx, obj, o.bindings, selfModifier); err != nil {
		return err
	} else {
		obj = merged
	}
	if err := internal.ApplyNamespacer(o.namespacer, &obj); err != nil {
		return err
	}
	internal.LogStart(logger, logging.Patch)
	return o.execute(ctx, obj)
}

func (o *operation) execute(ctx context.Context, obj unstructured.Unstructured) error {
	var lastErr error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		lastErr = o.tryPatchResource(ctx, obj)
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

func (o *operation) tryPatchResource(ctx context.Context, obj unstructured.Unstructured) error {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.ObjectKey(&obj), &actual)
	if err != nil {
		return err
	}
	return o.updateResource(ctx, &actual, obj)
}

func (o *operation) updateResource(ctx context.Context, actual *unstructured.Unstructured, obj unstructured.Unstructured) error {
	patched, err := client.PatchObject(actual, &obj)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(patched)
	if err != nil {
		return err
	}
	return o.handleCheck(ctx, obj, o.client.Patch(ctx, actual, ctrlclient.RawPatch(types.MergePatchType, bytes)))
}

func (o *operation) handleCheck(ctx context.Context, obj unstructured.Unstructured, err error) error {
	bindings := o.bindings
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
