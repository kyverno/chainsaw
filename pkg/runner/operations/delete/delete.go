package delete

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operation struct {
	client     client.Client
	obj        ctrlclient.Object
	namespacer namespacer.Namespacer
	check      *v1alpha1.Check
}

func New(client client.Client, obj ctrlclient.Object, namespacer namespacer.Namespacer, check *v1alpha1.Check) operations.Operation {
	return &operation{
		client:     client,
		obj:        obj,
		namespacer: namespacer,
		check:      check,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := logging.FromContext(ctx).WithResource(o.obj)
	defer func() {
		if err != nil {
			logger.Log(logging.Delete, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		} else {
			logger.Log(logging.Delete, logging.DoneStatus, color.BoldGreen)
		}
	}()
	if o.namespacer != nil {
		if err := o.namespacer.Apply(o.obj); err != nil {
			return err
		}
	}
	logger.Log(logging.Delete, logging.RunStatus, color.BoldFgCyan)
	return o.deleteResource(ctx, logger)
}

func (o *operation) deleteResource(ctx context.Context, logger logging.Logger) error {
	candidates, err := internal.Read(ctx, o.obj, o.client)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	var deleted []unstructured.Unstructured
	for i := range candidates {
		candidate := candidates[i]
		if err := o.tryDeleteCandidate(ctx, &candidate); err != nil {
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

func (o *operation) tryDeleteCandidate(ctx context.Context, candidate *unstructured.Unstructured) error {
	if err := o.client.Delete(ctx, candidate); err != nil && !kerrors.IsNotFound(err) {
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

func (o *operation) handleCheck(ctx context.Context, candidate *unstructured.Unstructured, err error) error {
	if o.check == nil || o.check.Value == nil {
		return err
	}
	bindings := binding.NewBindings()
	if err == nil {
		bindings.Register("$error", binding.NewBinding(nil))
	} else {
		bindings.Register("$error", binding.NewBinding(err.Error()))
	}
	errs, validationErr := assert.Validate(ctx, o.check.Value, candidate, bindings)
	if validationErr != nil {
		return validationErr
	}
	return errs.ToAggregate()
}
