package delete

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/check"
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
	expect     []v1alpha1.Expectation
}

func New(client client.Client, obj unstructured.Unstructured, namespacer namespacer.Namespacer, expect ...v1alpha1.Expectation) operations.Operation {
	return &operation{
		client:     client,
		obj:        obj,
		namespacer: namespacer,
		expect:     expect,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := internal.GetLogger(ctx, &o.obj)
	defer func() { internal.LogEnd(logger, logging.Delete, err) }()
	if err := internal.ApplyNamespacer(o.namespacer, &o.obj); err != nil {
		return err
	}
	internal.LogStart(logger, logging.Delete)
	return o.deleteResource(ctx)
}

func (o *operation) deleteResource(ctx context.Context) error {
	candidates, err := internal.Read(ctx, &o.obj, o.client)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	var deleted []unstructured.Unstructured
	for _, candidate := range candidates {
		if err := o.tryDeleteCandidate(ctx, candidate); err != nil {
			return err
		}
		deleted = append(deleted, candidate)
	}
	for i := range deleted {
		candidate := deleted[i]
		if err := o.waitForDeletion(ctx, &candidate); err != nil {
			return err
		}
	}
	return nil
}

func (o *operation) tryDeleteCandidate(ctx context.Context, candidate unstructured.Unstructured) error {
	if err := o.client.Delete(ctx, &candidate); err != nil && !kerrors.IsNotFound(err) {
		return o.handleCheck(ctx, candidate, err)
	}
	return o.handleCheck(ctx, candidate, nil)
}

func (o *operation) waitForDeletion(ctx context.Context, candidate *unstructured.Unstructured) error {
	gvk := candidate.GetObjectKind().GroupVersionKind()
	return wait.PollUntilContextCancel(ctx, internal.PollInterval, true, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(gvk)
		err := o.client.Get(ctx, client.ObjectKey(candidate), &actual)
		if kerrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}

func (o *operation) handleCheck(ctx context.Context, candidate unstructured.Unstructured, err error) error {
	bindings := binding.NewBindings()
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if matched, err := check.Expectations(ctx, candidate, bindings, o.expect...); matched {
		return err
	}
	return err
}
