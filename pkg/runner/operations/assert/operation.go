package assert

import (
	"context"
	"errors"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	operrors "github.com/kyverno/chainsaw/pkg/runner/operations/errors"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/chainsaw/pkg/runner/template"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
)

type operation struct {
	client     client.Client
	base       unstructured.Unstructured
	namespacer namespacer.Namespacer
	template   bool
}

func New(
	client client.Client,
	expected unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	template bool,
) operations.Operation {
	return &operation{
		client:     client,
		base:       expected,
		namespacer: namespacer,
		template:   template,
	}
}

func (o *operation) Exec(ctx context.Context, bindings binding.Bindings) (_ operations.Outputs, _err error) {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Assert, _err)
	}()
	if o.template {
		if err := template.ResourceRef(ctx, &obj, bindings); err != nil {
			return nil, err
		}
	}
	if obj.GetKind() != "" {
		if err := internal.ApplyNamespacer(o.namespacer, &obj); err != nil {
			return nil, err
		}
	}
	internal.LogStart(logger, logging.Assert)
	return nil, o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings binding.Bindings, obj unstructured.Unstructured) error {
	var lastErrs []error
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (_ bool, err error) {
		var errs []error
		defer func() {
			// record last errors only if there was no real error
			if err == nil {
				lastErrs = errs
			}
		}()
		if obj.GetAPIVersion() == "" || obj.GetKind() == "" {
			_errs, err := check.Check(ctx, nil, bindings, &v1alpha1.Check{Value: obj.UnstructuredContent()})
			if err != nil {
				return false, err
			}
			if len(_errs) != 0 {
				for _, _err := range _errs {
					errs = append(errs, _err)
				}
			} else {
				// at least one match found
				return true, nil
			}
		} else {
			if candidates, err := internal.Read(ctx, &obj, o.client); err != nil {
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
					_errs, err := check.Check(ctx, candidate.UnstructuredContent(), bindings, &v1alpha1.Check{Value: obj.UnstructuredContent()})
					if err != nil {
						return false, err
					}
					if len(_errs) != 0 {
						errs = append(errs, operrors.ResourceError(obj, candidate, _errs))
					} else {
						// at least one match found
						return true, nil
					}
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
