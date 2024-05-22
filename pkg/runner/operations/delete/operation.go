package delete

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/mutate"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operation struct {
	client            client.Client
	base              unstructured.Unstructured
	namespacer        namespacer.Namespacer
	template          bool
	expect            []v1alpha1.Expectation
	propagationPolicy metav1.DeletionPropagation
}

func New(
	client client.Client,
	obj unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	template bool,
	propagationPolicy metav1.DeletionPropagation,
	expect ...v1alpha1.Expectation,
) operations.Operation {
	return &operation{
		client:            client,
		base:              obj,
		namespacer:        namespacer,
		template:          template,
		expect:            expect,
		propagationPolicy: propagationPolicy,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (_ operations.Outputs, _err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Delete, _err)
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
	internal.LogStart(logger, logging.Delete)
	return nil, o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	resources, err := o.getResourcesToDelete(ctx, obj)
	if err != nil {
		return err
	}
	return o.deleteResources(ctx, bindings, resources...)
}

func (o *operation) getResourcesToDelete(ctx context.Context, obj unstructured.Unstructured) ([]unstructured.Unstructured, error) {
	resources, err := internal.Read(ctx, &obj, o.client)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return resources, nil
}

func (o *operation) deleteResources(ctx context.Context, bindings binding.Bindings, resources ...unstructured.Unstructured) error {
	var errs []error
	var deleted []unstructured.Unstructured
	for _, resource := range resources {
		err := o.deleteResource(ctx, resource)
		// if the resource was successfully deleted, record it to track actual deletion
		if err == nil {
			deleted = append(deleted, resource)
		}
		// check if the result was the expected one
		if err := o.handleCheck(ctx, bindings, resource, err); err != nil {
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
	if err := o.client.Delete(ctx, &resource, ctrlclient.PropagationPolicy(o.propagationPolicy)); err != nil {
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

func (o *operation) handleCheck(ctx context.Context, bindings binding.Bindings, resource unstructured.Unstructured, err error) error {
	if err == nil {
		bindings = apibindings.RegisterNamedBinding(ctx, bindings, "error", nil)
	} else {
		bindings = apibindings.RegisterNamedBinding(ctx, bindings, "error", err.Error())
	}
	if matched, err := check.Expectations(ctx, resource, bindings, o.expect...); matched {
		return err
	}
	return err
}
