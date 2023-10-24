package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Error(ctx context.Context, logger logging.Logger, expected ctrlclient.Object, c client.Client) (_err error) {
	attempts := 0
	defer func() {
		logging.ResourceOp(logger, "ERROR", client.ObjectKey(expected), expected, attempts, _err)
	}()
	return wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (bool, error) {
		attempts++
		candidates, err := read(ctx, expected, c)
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		if len(candidates) == 0 {
			return true, nil
		}
		return false, nil
	})
}
