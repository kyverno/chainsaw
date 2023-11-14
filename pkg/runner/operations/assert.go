package operations

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

type AssertOperation struct {
	baseOperation
	expected unstructured.Unstructured
}

func (a *AssertOperation) Name() string {
	return "ASSERT"
}

func (a *AssertOperation) Exec(ctx context.Context) (_err error) {
	const operation = "ASSERT"
	logger := logging.FromContext(ctx).WithResource(&a.expected)
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
		if candidates, err := read(ctx, &a.expected, a.client); err != nil {
			if kerrors.IsNotFound(err) {
				errs = append(errs, errors.New("actual resource not found"))
				return false, nil
			}
			return false, err
		} else if len(candidates) == 0 {
			errs = append(errs, errors.New("no actual resource found"))
		} else {
			for i := range candidates {
				candidate := candidates[i]
				_errs, err := assert.Validate(ctx, a.expected.UnstructuredContent(), candidate.UnstructuredContent(), nil)
				if err != nil {
					return false, err
				}
				if len(_errs) != 0 {
					var output []string
					for _, _err := range _errs {
						output = append(output, "    "+_err.Error())
					}
					slices.Sort(output)
					errs = append(errs, fmt.Errorf("resource %s doesn't match expectation:\n%s", client.Name(client.ObjectKey(&candidate)), strings.Join(output, "\n")))
				} else {
					// at least one match found
					return true, nil
				}
			}
		}
		return false, nil
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
