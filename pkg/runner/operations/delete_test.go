package operations

import (
	"context"
	"errors"
	"testing"
	"time"

	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_operationDelete(t *testing.T) {
	pod := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name": "test-pod",
			},
		},
	}
	tests := []struct {
		name         string
		object       ctrlclient.Object
		client       *tclient.FakeClient
		expectedErr  error
		expectedLogs []string
	}{
		{
			name:   "not found",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					return nil
				},
			},
			expectedErr:  nil,
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [DONE]"},
		},
		{
			name:   "failed get",
			object: pod.DeepCopy(),
			client: &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					return kerrors.NewInternalError(errors.New("failed to get the pod"))
				},
				DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					return nil
				},
			},
			expectedErr:  kerrors.NewInternalError(errors.New("failed to get the pod")),
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [ERROR\nInternal error occurred: failed to get the pod]"},
		},
		{
			name:   "failed delete",
			object: pod.DeepCopy(),
			client: &tclient.FakeClient{
				GetFn: func(_ context.Context, _ int, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					return nil
				},
				DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					return kerrors.NewInternalError(errors.New("failed to delete the pod"))
				},
			},
			expectedErr:  kerrors.NewInternalError(errors.New("failed to delete the pod")),
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [ERROR\nInternal error occurred: failed to delete the pod]"},
		},
		{
			name:   "ok",
			object: pod.DeepCopy(),
			client: &tclient.FakeClient{
				GetFn: func(_ context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					if call < 10 {
						return nil
					}
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					return nil
				},
			},
			expectedErr:  nil,
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [DONE]"},
		},
		{
			name:   "poll succeeds but returns error after",
			object: pod.DeepCopy(),
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					return nil
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
					return nil
				},
			},
			expectedErr:  context.DeadlineExceeded,
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [ERROR\ncontext deadline exceeded]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			ctx := logging.IntoContext(ctxt, logger)
			err := operationDelete(ctx, tt.object, tt.client)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
