package operations

import (
	"context"
	"fmt"
	"testing"

	fakeClient "github.com/kyverno/chainsaw/pkg/runner/client"
	fakeLogger "github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		object       ctrlclient.Object
		client       *fakeClient.FakeClient
		expectedErr  error
		expectedLogs []string
	}{
		{
			name: "Successful delete",
			object: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name": "test-pod",
					},
				},
			},
			client: &fakeClient.FakeClient{
				GetFake: func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return errors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				DeleteFake: func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					return nil
				},
			},
			expectedErr:  nil,
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [DONE]"},
		},
		{
			name: "Failed delete",
			object: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name": "bad-test-pod",
					},
				},
			},
			client: &fakeClient.FakeClient{
				GetFake: func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
				DeleteFake: func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					return errors.NewInternalError(fmt.Errorf("failed to delete the pod"))
				},
			},
			expectedErr:  errors.NewInternalError(fmt.Errorf("failed to delete the pod")),
			expectedLogs: []string{"DELETE: [RUNNING...]", "DELETE: [ERROR\nInternal error occurred: failed to delete the pod]"},
		},
		///

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &fakeLogger.MockLogger{}

			err := Delete(context.TODO(), logger, tt.object, tt.client)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
