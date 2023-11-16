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
	dryRun  bool
	cleaner cleanup.Cleaner
	check   interface{}
	created bool
}

func (a *operation) Cleanup() {
	if a.cleaner != nil && a.created && !a.dryRun {
		a.cleaner(a.obj, a.client)
	}
}

func New(client client.Client, obj ctrlclient.Object, dryRun bool, cleaner cleanup.Cleaner, check interface{}) *operation {
	return &operation{
		client:  client,
		obj:     obj,
		dryRun:  dryRun,
		cleaner: cleaner,
		check:   check,
	}
}

func (a *operation) Exec(ctx context.Context) (_err error) {
	const operation = "APPLY"
	logger := logging.FromContext(ctx).WithResource(a.obj)
	logger.Log(operation, color.BoldFgCyan, "RUNNING...")
	defer func() {
		if _err == nil {
			logger.Log(operation, color.BoldGreen, "DONE")
		} else {
			logger.Log(operation, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	return wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(a.obj.GetObjectKind().GroupVersionKind())
		err := a.client.Get(ctx, client.ObjectKey(a.obj), &actual)
		if err == nil {
			patchOptions := []ctrlclient.PatchOption{}
			if a.dryRun {
				patchOptions = append(patchOptions, ctrlclient.DryRunAll)
			}
			patched, err := client.PatchObject(&actual, a.obj)
			if err != nil {
				return false, err
			}
			bytes, err := json.Marshal(patched)
			if err != nil {
				return false, err
			}
			err = a.client.Patch(ctx, &actual, ctrlclient.RawPatch(types.MergePatchType, bytes), patchOptions...)
			if a.check == nil {
				return err == nil, err
			} else {
				actual := map[string]interface{}{
					"error":    nil,
					"resource": a.obj,
				}
				if err != nil {
					actual["error"] = err.Error()
				}
				errs, err := assert.Validate(ctx, a.check, actual, nil)
				if err != nil {
					return false, err
				}
				return true, errs.ToAggregate()
			}
		} else if kerrors.IsNotFound(err) {
			var createOptions []ctrlclient.CreateOption
			if a.dryRun {
				createOptions = append(createOptions, ctrlclient.DryRunAll)
			}
			err := a.client.Create(ctx, a.obj, createOptions...)
			if a.check == nil {
				return err == nil, err
			} else {
				actual := map[string]interface{}{
					"error":    nil,
					"resource": a.obj,
				}
				if err != nil {
					actual["error"] = err.Error()
				}
				errs, err := assert.Validate(ctx, a.check, actual, nil)
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
