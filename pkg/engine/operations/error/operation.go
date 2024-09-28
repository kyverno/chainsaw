package error

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/checks"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/engine/templating"
	"go.uber.org/multierr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
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

func (o *operation) Exec(ctx context.Context, bindings apis.Bindings) (_ outputs.Outputs, _err error) {
	if bindings == nil {
		bindings = apis.NewBindings()
	}
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Error, _err)
	}()
	if o.template {
		if err := templating.ResourceRef(ctx, &obj, bindings); err != nil {
			return nil, err
		}
	}
	if obj.GetKind() != "" {
		if err := internal.ApplyNamespacer(o.namespacer, o.client, &obj); err != nil {
			return nil, err
		}
	}
	internal.LogStart(logger, logging.Error)
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
			_errs, err := checks.Check(ctx, nil, bindings, ptr.To(v1alpha1.NewCheck(obj.UnstructuredContent())))
			if err != nil {
				return false, err
			}
			if len(_errs) == 0 {
				errs = append(errs, fmt.Errorf("expectation matched"))
			}
			return len(errs) == 0, nil
		} else {
			if candidates, err := internal.Read(ctx, &obj, o.client); err != nil {
				if kerrors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			} else if len(candidates) == 0 {
				return true, nil
			} else {
				for i := range candidates {
					candidate := candidates[i]
					_errs, err := checks.Check(ctx, candidate.UnstructuredContent(), bindings, ptr.To(v1alpha1.NewCheck(obj.UnstructuredContent())))
					if err != nil {
						return false, err
					}
					if len(_errs) == 0 {
						errs = append(errs, fmt.Errorf("%s/%s/%s - resource matches expectation", candidate.GetAPIVersion(), candidate.GetKind(), client.Name(client.Key(&candidate))))
					}
				}
				return len(errs) == 0, nil
			}
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
