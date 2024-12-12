package assert

import (
	"context"
	"errors"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	operrors "github.com/kyverno/chainsaw/pkg/engine/operations/errors"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
)

type operation struct {
	compilers  compilers.Compilers
	client     client.Client
	base       unstructured.Unstructured
	namespacer namespacer.Namespacer
	template   bool
}

func New(
	compilers compilers.Compilers,
	client client.Client,
	expected unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	template bool,
) operations.Operation {
	return &operation{
		compilers:  compilers,
		client:     client,
		base:       expected,
		namespacer: namespacer,
		template:   template,
	}
}

func (o *operation) Exec(ctx context.Context, bindings apis.Bindings) (_ outputs.Outputs, _err error) {
	if bindings == nil {
		bindings = apis.NewBindings()
	}
	obj := o.base
	defer func() {
		internal.LogEnd(ctx, logging.Assert, &obj, _err)
	}()
	if o.template {
		if err := templating.ResourceRef(ctx, o.compilers, &obj, bindings); err != nil {
			return nil, err
		}
	}
	if obj.GetKind() != "" {
		if err := internal.ApplyNamespacer(o.namespacer, o.client, &obj); err != nil {
			return nil, err
		}
	}
	internal.LogStart(ctx, logging.Assert, &obj)
	return nil, o.execute(ctx, bindings, obj)
}

func (o *operation) execute(ctx context.Context, bindings apis.Bindings, obj unstructured.Unstructured) error {
	var lastErrs []error
	err := wait.PollUntilContextCancel(ctx, client.PollInterval, false, func(ctx context.Context) (_ bool, err error) {
		var errs []error
		defer func() {
			// record last errors only if there was no real error
			if err == nil {
				lastErrs = errs
			}
		}()
		if obj.GetAPIVersion() == "" || obj.GetKind() == "" {
			_errs, err := checks.Check(ctx, o.compilers, nil, bindings, ptr.To(v1alpha1.NewCheck(obj.UnstructuredContent())))
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
					_errs, err := checks.Check(ctx, o.compilers, candidate.UnstructuredContent(), bindings, ptr.To(v1alpha1.NewCheck(obj.UnstructuredContent())))
					if err != nil {
						return false, err
					}
					if len(_errs) != 0 {
						errs = append(errs, operrors.ResourceError(o.compilers, obj, candidate, o.template, bindings, _errs))
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