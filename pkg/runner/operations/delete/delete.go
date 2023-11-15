package delete

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type operation struct {
	client client.Client
	obj    ctrlclient.Object
}

func New(client client.Client, obj ctrlclient.Object) *operation {
	return &operation{
		client: client,
		obj:    obj,
	}
}

func (d *operation) Exec(ctx context.Context) (_err error) {
	const operation = "DELETE"
	logger := logging.FromContext(ctx).WithResource(d.obj)
	logger.Log(operation, color.BoldFgCyan, "RUNNING...")
	defer func() {
		if _err == nil {
			logger.Log(operation, color.BoldGreen, "DONE")
		} else {
			logger.Log(operation, color.BoldRed, fmt.Sprintf("ERROR\n%s", _err))
		}
	}()
	candidates, _err := internal.Read(ctx, d.obj, d.client)
	if _err != nil {
		if errors.IsNotFound(_err) {
			return nil
		}
		return _err
	}
	for i := range candidates {
		err := d.client.Delete(ctx, &candidates[i])
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}
	gvk := d.obj.GetObjectKind().GroupVersionKind()
	for i := range candidates {
		if err := wait.PollUntilContextCancel(ctx, internal.PollInterval, true, func(ctx context.Context) (bool, error) {
			var actual unstructured.Unstructured
			actual.SetGroupVersionKind(gvk)
			err := d.client.Get(ctx, client.ObjectKey(&candidates[i]), &actual)
			if err != nil {
				if errors.IsNotFound(err) {
					return true, nil
				}
				return false, err
			}
			return false, nil
		}); err != nil {
			return err
		}
	}
	return nil
}
