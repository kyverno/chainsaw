package assert

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno-json/pkg/engine/assert"
	"github.com/kyverno/kyverno/ext/output/color"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
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
	return o.pollForAssertion(ctx)
}

func (o *operation) pollForAssertion(ctx context.Context) error {
	var lastErrs []error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (bool, error) {
		candidates, errs, updateErrs, err := o.fetchAndValidateCandidates(ctx)
		if err != nil || len(errs) != 0 {
			if updateErrs {
				lastErrs = errs
			}
			return false, err
		}
		return len(candidates) > 0, nil
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

func (o *operation) fetchAndValidateCandidates(ctx context.Context) ([]unstructured.Unstructured, []error, bool, error) {
	candidates, err := internal.Read(ctx, &o.expected, o.client)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, []error{errors.New("actual resource not found")}, true, nil
		}
		return nil, nil, false, err
	}

	if len(candidates) == 0 {
		return nil, []error{errors.New("no actual resource found")}, true, nil
	}

	var errs []error
	for _, candidate := range candidates {
		if _errs, err := assert.Validate(ctx, o.expected.UnstructuredContent(), candidate.UnstructuredContent(), nil); err != nil {
			return nil, nil, false, err
		} else if len(_errs) != 0 {
			errs = appendErrorDetails(errs, candidate, _errs)
		}
	}
	return candidates, errs, true, nil
}

func appendErrorDetails(errs []error, candidate unstructured.Unstructured, _errs field.ErrorList) []error {
	for _, _err := range _errs {
		errs = append(errs, fmt.Errorf("%s/%s/%s - %w", candidate.GetAPIVersion(), candidate.GetKind(), client.Name(client.ObjectKey(&candidate)), _err))
	}
	return errs
}
