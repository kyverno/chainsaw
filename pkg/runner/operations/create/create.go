package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func (c *operation) Cleanup() {
	if c.cleaner != nil && c.created && !c.dryRun {
		c.cleaner(c.obj, c.client)
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

func (c *operation) Exec(ctx context.Context) (_err error) {
	const operation = "CREATE"
	logger := logging.FromContext(ctx).WithResource(c.obj)
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
		actual.SetGroupVersionKind(c.obj.GetObjectKind().GroupVersionKind())
		err := c.client.Get(ctx, client.ObjectKey(c.obj), &actual)
		if err == nil {
			return false, errors.New("the resource already exists in the cluster")
		} else if kerrors.IsNotFound(err) {
			var createOptions []ctrlclient.CreateOption
			if c.dryRun {
				createOptions = append(createOptions, ctrlclient.DryRunAll)
			}
			err := c.client.Create(ctx, c.obj, createOptions...)
			if c.check == nil {
				return err == nil, err
			} else {
				actual := map[string]interface{}{
					"error":    nil,
					"resource": c.obj,
				}
				if err != nil {
					actual["error"] = err.Error()
				}
				errs, err := assert.Validate(ctx, c.check, actual, nil)
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
