package create

import (
	"context"
	"errors"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	kjsonv1alpha1 "github.com/kyverno/kyverno-json/pkg/apis/v1alpha1"
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
	cleaner    cleanup.Cleaner
	check      *kjsonv1alpha1.Any
}

func New(client client.Client, obj ctrlclient.Object, namespacer namespacer.Namespacer, cleaner cleanup.Cleaner, check *kjsonv1alpha1.Any) operations.Operation {
	return &operation{
		client:     client,
		obj:        obj,
		namespacer: namespacer,
		cleaner:    cleaner,
		check:      check,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := logging.FromContext(ctx).WithResource(o.obj)
	defer func() {
		if err != nil {
			logger.Log(logging.Create, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		} else {
			logger.Log(logging.Create, logging.DoneStatus, color.BoldGreen)
		}
	}()
	if o.namespacer != nil {
		if err := o.namespacer.Apply(o.obj); err != nil {
			return err
		}
	}
	logger.Log(logging.Create, logging.RunStatus, color.BoldFgCyan)
	return o.createResource(ctx, logger)
}

func (o *operation) createResource(ctx context.Context, logger logging.Logger) error {
	return wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		err := o.tryCreateResource(ctx)
		return err == nil, err
	})
}

func (o *operation) tryCreateResource(ctx context.Context) error {
	var actual unstructured.Unstructured
	actual.SetGroupVersionKind(o.obj.GetObjectKind().GroupVersionKind())
	err := o.client.Get(ctx, client.ObjectKey(o.obj), &actual)
	if err == nil {
		return errors.New("the resource already exists in the cluster")
	}
	if kerrors.IsNotFound(err) {
		return o.create_Resource(ctx)
	}
	return err
}

func (o *operation) create_Resource(ctx context.Context) error {
	err := o.client.Create(ctx, o.obj)
	if err == nil && o.cleaner != nil {
		o.cleaner(o.obj, o.client)
	}
	return o.handleCheck(ctx, err)
}

func (o *operation) handleCheck(ctx context.Context, err error) error {
	if o.check == nil || o.check.Value == nil {
		return err
	}
	actual := map[string]interface{}{
		"error":    nil,
		"resource": o.obj,
	}
	if err != nil {
		actual["error"] = err.Error()
	}
	errs, validationErr := assert.Validate(ctx, o.check.Value, actual, nil)
	if validationErr != nil {
		return validationErr
	}
	return errs.ToAggregate()
}
