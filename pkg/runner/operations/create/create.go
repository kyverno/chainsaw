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
	check      interface{}
}

func New(client client.Client, obj ctrlclient.Object, namespacer namespacer.Namespacer, cleaner cleanup.Cleaner, check interface{}) operations.Operation {
	return &operation{
		client:     client,
		obj:        obj,
		namespacer: namespacer,
		cleaner:    cleaner,
		check:      check,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	logger := logging.FromContext(ctx).WithResource(o.obj)
	defer func() {
		if _err == nil {
			logger.Log(logging.Create, logging.DoneStatus, color.BoldGreen)
		} else {
			logger.Log(logging.Create, logging.ErrorStatus, color.BoldRed, logging.ErrSection(_err))
		}
	}()
	if o.namespacer != nil {
		if err := o.namespacer.Apply(o.obj); err != nil {
			return err
		}
	}
	logger.Log(logging.Create, logging.RunStatus, color.BoldFgCyan)
	return wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(o.obj.GetObjectKind().GroupVersionKind())
		err := o.client.Get(ctx, client.ObjectKey(o.obj), &actual)
		if err == nil {
			return false, errors.New("the resource already exists in the cluster")
		} else if kerrors.IsNotFound(err) {
			err := o.client.Create(ctx, o.obj)
			if err == nil && o.cleaner != nil {
				o.cleaner(o.obj, o.client)
			}
			if o.check == nil {
				return err == nil, err
			} else {
				actual := map[string]interface{}{
					"error":    nil,
					"resource": o.obj,
				}
				if err != nil {
					actual["error"] = err.Error()
				}
				errs, err := assert.Validate(ctx, o.check, actual, nil)
				if err != nil {
					return false, err
				}
				return true, errs.ToAggregate()
			}
		} else {
			return false, err
		}
	})
}
