package apply

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operation struct {
	client  client.Client
	obj     ctrlclient.Object
	cleaner cleanup.Cleaner
	check   interface{}
}

func New(client client.Client, obj ctrlclient.Object, cleaner cleanup.Cleaner, check interface{}) *operation {
	return &operation{
		client:  client,
		obj:     obj,
		cleaner: cleaner,
		check:   check,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	const operationName = "APPLY"
	logger := logging.FromContext(ctx).WithResource(o.obj)

	logger.Log(operationName, color.BoldFgCyan, "RUNNING...")
	defer func() {
		if _err == nil {
			logger.Log(operationName, color.BoldGreen, "DONE")
		} else {
			logger.Log(operationName, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()

	return wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(o.obj.GetObjectKind().GroupVersionKind())
		err := o.client.Get(ctx, client.ObjectKey(o.obj), &actual)
		if err == nil {
			return o.patchResource(ctx, &actual)
		} else if kerrors.IsNotFound(err) {
			return o.createResource(ctx)
		} else {
			return false, err
		}
	})
}

func (o *operation) createResource(ctx context.Context) (bool, error) {
	if err := o.client.Create(ctx, o.obj); err != nil {
		if o.check != nil {
			return o.performCheck(ctx, err)
		}
		return false, err
	}
	if o.check != nil {
		return o.performCheck(ctx, nil)
	}
	if o.cleaner != nil {
		o.cleaner(o.obj, o.client)
	}
	return true, nil
}

func (o *operation) patchResource(ctx context.Context, actual *unstructured.Unstructured) (bool, error) {
	patched, err := client.PatchObject(actual, o.obj)
	if err != nil {
		return false, err
	}
	bytes, err := json.Marshal(patched)
	if err != nil {
		return false, err
	}
	if err := o.client.Patch(ctx, actual, ctrlclient.RawPatch(types.MergePatchType, bytes)); err != nil {
		if o.check != nil {
			return o.performCheck(ctx, err)
		}
		return false, err
	}
	if o.check != nil {
		return o.performCheck(ctx, nil)
	}
	return o.updateStatus(ctx)
}

func (o *operation) performCheck(ctx context.Context, err error) (bool, error) {
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

func (o *operation) updateStatus(ctx context.Context) (bool, error) {
	if err := o.client.Status().Update(ctx, o.obj); err != nil {
		if o.check != nil {
			return o.performCheck(ctx, err)
		}
		return false, err
	}
	if o.check != nil {
		return o.performCheck(ctx, nil)
	}
	return true, nil
}
