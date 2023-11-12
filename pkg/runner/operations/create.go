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

func operationCreate(ctx context.Context, obj ctrlclient.Object, c client.Client, shouldFail bool, dryRun bool, cleaner cleanup.Cleaner) (_err error) {
	const operation = "CREATE"
	logger := logging.FromContext(ctx).WithResource(obj)

	if dryRun {
		logger.Log(operation, color.BoldFgCyan, "DRY RUN INITIATED")

		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
		err := c.Get(ctx, client.ObjectKey(obj), &actual)
		if err == nil {
			logger.Log(operation, color.BoldYellow, "DRY RUN: Resource already exists")
		} else if kerrors.IsNotFound(err) {
			logger.Log(operation, color.BoldYellow, "DRY RUN: Resource can be created")
		} else {
			logger.Log(operation, color.BoldYellow, fmt.Sprintf("DRY RUN: Error checking resource existence: %s", err))
		}

		logger.Log(operation, color.BoldFgCyan, "DRY RUN COMPLETED")
		return nil
	}

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
		actual.SetGroupVersionKind(obj.GetObjectKind().GroupVersionKind())
		err := c.Get(ctx, client.ObjectKey(obj), &actual)
		if err == nil {
			return false, errors.New("the resource already exists in the cluster")
		} else if kerrors.IsNotFound(err) {
			if err := c.Create(ctx, obj); err != nil {
				if shouldFail {
					return true, nil
				}
				return false, err
			} else {
				if cleaner != nil {
					cleaner(obj, c)
				}
				if shouldFail {
					return false, errors.New("an error was expected but didn't happen")
				}
			}
		} else {
			return false, err
		}
		return true, nil
	})
}
