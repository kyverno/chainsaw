package operations

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func operationError(ctx context.Context, expected unstructured.Unstructured, c client.Client) (_err error) {
	const operation = "ERROR "
	logger := logging.FromContext(ctx).WithResource(&expected)
	logger.Log(operation, color.BoldFgCyan, "RUNNING...")
	defer func() {
		if _err == nil {
			logger.Log(operation, color.BoldGreen, "DONE")
		} else {
			logger.Log(operation, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	var lastErrs []error
	err := wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (_ bool, err error) {
		var errs []error
		defer func() {
			// record last errors only if there was no real error
			if err == nil {
				lastErrs = errs
			}
		}()
		if candidates, err := read(ctx, &expected, c); err != nil {
			if kerrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		} else if len(candidates) == 0 {
			return true, nil
		} else {
			for i := range candidates {
				candidate := candidates[i]
				_errs, err := assert.Validate(ctx, expected.UnstructuredContent(), candidate.UnstructuredContent(), nil)
				if err != nil {
					return false, err
				}
				if len(_errs) == 0 {
					for _, _err := range _errs {
						errs = append(errs, _err)
					}
				}
				// if err := match.Match(expected.UnstructuredContent(), candidate.UnstructuredContent()); err == nil {
				// 	// at least one match found
				// 	errs = append(errs, fmt.Errorf("found an actual resource matching expectation (%s/%s / %s)", candidate.GetAPIVersion(), candidate.GetKind(), client.ObjectKey(&candidate)))
				// }
			}
			return len(errs) == 0, nil
		}
	})
	// if no error, return success
	if err == nil {
		return nil
	}
	// eventually return a combination of last errors
	if len(lastErrs) != 0 {
		return multierr.Combine(lastErrs...)
	}
	// return received error
	return err
}
