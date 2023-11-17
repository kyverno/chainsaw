package delete

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno/ext/output/color"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operation struct {
	client     client.Client
	obj        ctrlclient.Object
	namespacer namespacer.Namespacer
}

func New(client client.Client, obj ctrlclient.Object, namespacer namespacer.Namespacer) operations.Operation {
	return &operation{
		client:     client,
		obj:        obj,
		namespacer: namespacer,
	}
}

func (o *operation) Exec(ctx context.Context) error {
	logger := logging.FromContext(ctx).WithResource(o.obj)

	if o.namespacer != nil {
		if err := o.namespacer.Apply(o.obj); err != nil {
			return err
		}
	}

	logger.Log(logging.Delete, logging.RunStatus, color.BoldFgCyan)

	return o.deleteResource(ctx, logger)
}

func (o *operation) deleteResource(ctx context.Context, logger logging.Logger) error {
	candidates, err := internal.Read(ctx, o.obj, o.client)
	if err != nil {
		if kerrors.IsNotFound(err) {
			logger.Log(logging.Delete, logging.DoneStatus, color.BoldGreen)
			return nil
		}
		logger.Log(logging.Delete, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		return err
	}

	for i := range candidates {
		candidate := candidates[i]
		if err := o.tryDeleteCandidate(ctx, &candidate); err != nil {
			logger.Log(logging.Delete, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			return err
		}
		if err := o.waitForDeletion(ctx, &candidate); err != nil {
			logger.Log(logging.Delete, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
			return err
		}
	}

	logger.Log(logging.Delete, logging.DoneStatus, color.BoldGreen)
	return nil
}

func (o *operation) tryDeleteCandidate(ctx context.Context, candidate *unstructured.Unstructured) error {
	err := o.client.Delete(ctx, candidate)
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (o *operation) waitForDeletion(ctx context.Context, candidate *unstructured.Unstructured) error {
	gvk := candidate.GetObjectKind().GroupVersionKind()
	return wait.PollUntilContextCancel(ctx, internal.PollInterval, true, func(ctx context.Context) (bool, error) {
		var actual unstructured.Unstructured
		actual.SetGroupVersionKind(gvk)
		err := o.client.Get(ctx, client.ObjectKey(candidate), &actual)
		if kerrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}
