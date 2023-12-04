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
	"go.uber.org/multierr"
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
	defer func() {
		internal.LogEnd(logger, logging.Delete, err)
	}()
	if err := internal.ApplyNamespacer(o.namespacer, &o.obj); err != nil {
		return err
	}
	internal.LogStart(logger, logging.Delete)
	return o.execute(ctx)
}

func (o *operation) execute(ctx context.Context) error {
	resources, err := o.getResourcesToDelete(ctx)
	if err != nil {
		return err
	}
	return o.deleteResources(ctx, resources...)
}

func (o *operation) getResourcesToDelete(ctx context.Context) ([]unstructured.Unstructured, error) {
	resources, err := internal.Read(ctx, &o.obj, o.client)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return resources, nil
}

func (o *operation) deleteResources(ctx context.Context, resources ...unstructured.Unstructured) error {
	var errs []error
	var deleted []unstructured.Unstructured
	for _, resource := range resources {
		err := o.deleteResource(ctx, resource)
		// if the resource was successfully deleted, record it to track actual deletion
		if err == nil {
			deleted = append(deleted, resource)
		}
		// check if the result was the expected one
		if err := o.handleCheck(ctx, resource, err); err != nil {
			errs = append(errs, err)
		}
	}
	for _, resource := range deleted {
		if err := o.waitForDeletion(ctx, resource); err != nil {
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}

func (o *operation) deleteResource(ctx context.Context, resource unstructured.Unstructured) error {
	if err := o.client.Delete(ctx, &resource); err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

func (o *operation) waitForDeletion(ctx context.Context, resource unstructured.Unstructured) error {
	gvk := resource.GetObjectKind().GroupVersionKind()
	key := client.ObjectKey(&resource)
	return wait.PollUntilContextCancel(ctx, internal.PollInterval, true, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(gvk)
		if err := o.client.Get(ctx, key, &actual); err != nil {
			if kerrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}

func (o *operation) handleCheck(ctx context.Context, resource unstructured.Unstructured, err error) error {
	bindings := binding.NewBindings()
	if err == nil {
		bindings = bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings = bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	if matched, err := check.Expectations(ctx, resource, bindings, o.expect...); matched {
		return err
	}
	return err
}
