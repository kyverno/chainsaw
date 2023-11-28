package assert

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno/ext/output/color"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

type operation struct {
	client     client.Client
	expected   unstructured.Unstructured
	namespacer namespacer.Namespacer
}

func New(client client.Client, expected unstructured.Unstructured, namespacer namespacer.Namespacer) operations.Operation {
	return &operation{
		client:     client,
		expected:   expected,
		namespacer: namespacer,
	}
}

func (o *operation) Exec(ctx context.Context) (_err error) {
	logger := logging.FromContext(ctx).WithResource(&o.expected)
	defer func() {
		if _err == nil {
			logger.Log(logging.Assert, logging.DoneStatus, color.BoldGreen)
		} else {
			logger.Log(logging.Assert, logging.ErrorStatus, color.BoldRed, logging.ErrSection(_err))
		}
	}()
	if o.namespacer != nil {
		if err := o.namespacer.Apply(&o.expected); err != nil {
			return err
		}
	}
	logger.Log(logging.Assert, logging.RunStatus, color.BoldFgCyan)
	var lastErrs []error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (_ bool, err error) {
		var errs []error
		defer func() {
			// record last errors only if there was no real error
			if err == nil {
				lastErrs = errs
			}
		}()
		if candidates, err := internal.Read(ctx, &o.expected, o.client); err != nil {
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
				_errs, err := check.Check(ctx, candidate.UnstructuredContent(), nil, &v1alpha1.Check{Value: o.expected.UnstructuredContent()})
				if err != nil {
					return false, err
				}
				if len(_errs) != 0 {
					for _, _err := range _errs {
						errs = append(errs, fmt.Errorf("%s/%s/%s - %w", candidate.GetAPIVersion(), candidate.GetKind(), client.Name(client.ObjectKey(&candidate)), _err))
					}
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
