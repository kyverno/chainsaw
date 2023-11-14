package operations

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type CreateOperation struct {
	baseOperation
	obj        ctrlclient.Object
	dryRun     bool
	cleaner    cleanup.Cleaner
	shouldFail bool
}

func (c *CreateOperation) Name() string {
	return "CREATE"
}

func (c *CreateOperation) Exec(ctx context.Context) (_err error) {
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

	return wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
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
			if err := c.client.Create(ctx, c.obj, createOptions...); err != nil {
				if c.shouldFail {
					return true, nil
				}
				return false, err
			} else {
				if c.cleaner != nil && !c.dryRun {
					c.cleaner(c.obj, c.client)
				}
				if c.shouldFail {
					return false, errors.New("an error was expected but didn't happen")
				}
				if c.dryRun {
					logger.Log(operation, color.BoldYellow, "DRY RUN: Resource creation simulated")
				}
			}
		} else {
			return false, err
		}
		return true, nil
	})
}
