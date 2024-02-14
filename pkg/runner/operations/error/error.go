package error

import (
	"context"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/check"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
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
	bindings   binding.Bindings
	template   bool
}

func New(
	client client.Client,
	expected unstructured.Unstructured,
	namespacer namespacer.Namespacer,
	bindings binding.Bindings,
	template bool,
) operations.Operation {
	if bindings == nil {
		bindings = binding.NewBindings()
	}
	return &operation{
		client:     client,
		base:       expected,
		namespacer: namespacer,
		bindings:   bindings,
		template:   template,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	obj := o.base
	logger := internal.GetLogger(ctx, &obj)
	defer func() {
		internal.LogEnd(logger, logging.Error, err)
	}()
	if o.template {
		if err := template.ResourceRef(ctx, &obj, o.bindings); err != nil {
			return err
		}
	}
	if obj.GetKind() != "" {
		if err := internal.ApplyNamespacer(o.namespacer, &obj); err != nil {
			return err
		}
	}
	internal.LogStart(logger, logging.Error)
	return o.execute(ctx, obj)
}

func (o *operation) execute(ctx context.Context, obj unstructured.Unstructured) error {
	var lastErrs []error
	bindings := o.bindings
	err := wait.PollUntilContextCancel(ctx, internal.PollInterval, false, func(ctx context.Context) (_ bool, err error) {
		var errs []error
		defer func() {
			// record last errors only if there was no real error
			if err == nil {
				lastErrs = errs
			}
		}()
		if !(obj.GetAPIVersion() != "" && obj.GetKind() != "") {
			_errs, err := check.Check(ctx, nil, bindings, &v1alpha1.Check{Value: obj.UnstructuredContent()})
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
					_errs, err := check.Check(ctx, candidate.UnstructuredContent(), bindings, &v1alpha1.Check{Value: obj.UnstructuredContent()})
					if err != nil {
						return false, err
					}
					if len(_errs) == 0 {
						errs = append(errs, fmt.Errorf("%s/%s/%s - resource matches expectation", candidate.GetAPIVersion(), candidate.GetKind(), client.Name(client.ObjectKey(&candidate))))
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
