package operations

import (
	"context"
	"errors"
	"testing"

	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_apply(t *testing.T) {
	podv1 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name": "test-pod",
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "test-container",
						"image": "test-image:v1",
					},
				},
			},
		},
	}
	podv2 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name": "test-pod",
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  "test-container",
						"image": "test-image:v2",
					},
				},
			},
		},
	}
	tests := []struct {
		name        string
		object      ctrlclient.Object
		client      *tclient.FakeClient
		shouldFail  bool
		expectedErr error
	}{{
		name:   "Resource already exists, patch it",
		object: podv2.DeepCopy(),
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = *podv1.DeepCopy()
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		shouldFail:  false,
		expectedErr: nil,
	}, {
		name:   "Resource does not exist, create it",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		shouldFail:  false,
		expectedErr: nil,
	}, {
		name:   "Error while getting resource",
		object: podv1.DeepCopy(),
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return errors.New("some arbitrary error")
			},
		},
		shouldFail:  true,
		expectedErr: errors.New("some arbitrary error"),
	},
		{
			name:   "Fail to patch existing resource",
			object: podv2.DeepCopy(),
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					*obj.(*unstructured.Unstructured) = *podv1.DeepCopy()
					return nil
				},
				PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
					return errors.New("patch failed")
				},
			},
			shouldFail:  false,
			expectedErr: errors.New("patch failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctx := logging.IntoContext(context.TODO(), logger)
			err := operationApply(ctx, tt.object, tt.client, tt.shouldFail, nil)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
