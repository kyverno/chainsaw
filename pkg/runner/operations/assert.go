package operations

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/match"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Assert(ctx context.Context, expected unstructured.Unstructured, c client.Client) error {
	var lastErrs []error
	err := wait.PollUntilContextCancel(ctx, interval, false, func(ctx context.Context) (_ bool, err error) {
		var errs []error
		defer func() {
			// record last errors only if there was no real error
			if err == nil {
				lastErrs = errs
			}
		}()
		if candidates, err := read(ctx, expected, c); err != nil {
			if kerrors.IsNotFound(err) {
				errs = append(errs, errors.New("actual resource not found"))
				return false, nil
			}
			return false, err
		} else if len(candidates) == 0 {
			errs = append(errs, errors.New("no actual resource found"))
		} else {
			for _, candidate := range candidates {
				if err := match.Match(expected.UnstructuredContent(), candidate.UnstructuredContent()); err != nil {
					diffStr, err := diff(expected, candidate)
					if err != nil {
						return false, err
					}
					errs = append(errs, fmt.Errorf("actual resource doesn't match expectation\n%s", diffStr))
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
