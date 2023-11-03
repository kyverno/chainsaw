package operations

import (
	"context"
	"errors"
	"testing"
	"time"

	fakeClient "github.com/kyverno/chainsaw/pkg/runner/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	fakeLogger "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clock "k8s.io/utils/clock/testing"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_apply(t *testing.T) {
	tests := []struct {
		name         string
		initialState *unstructured.Unstructured
		object       ctrlclient.Object
		client       *fakeClient.FakeClient
		shouldFail   bool
		expectedErr  error
	}{{
		name: "Resource already exists, patch it",
		initialState: &unstructured.Unstructured{
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
		},
		object: &unstructured.Unstructured{
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
		},
		client:      &fakeClient.FakeClient{},
		shouldFail:  false,
		expectedErr: nil,
	}, {
		name:         "Resource does not exist, create it",
		initialState: nil,
		object: &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name": "new-pod",
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
		},
		client:      &fakeClient.FakeClient{},
		shouldFail:  false,
		expectedErr: nil,
	}, {
		name:         "Error while getting resource",
		initialState: nil,
		object: &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name": "new-pod",
				},
			},
		},
		client:      &fakeClient.FakeClient{},
		shouldFail:  true,
		expectedErr: errors.New("some arbitrary error"),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClock := clock.NewFakePassiveClock(time.Now())
			ctx := logging.IntoContext(context.TODO(), logging.NewLogger(&fakeLogger.FakeLogger{}, fakeClock, "testName", "stepName"))
			tt.client.T = t
			tt.client.GetFake = func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				t.Helper()
				if tt.shouldFail {
					return tt.expectedErr
				}
				if tt.initialState != nil {
					*obj.(*unstructured.Unstructured) = *tt.initialState
				}
				return nil
			}
			tt.client.PatchFake = func(ctx context.Context, t *testing.T, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
				t.Helper()
				if tt.initialState != nil {
					*obj.(*unstructured.Unstructured) = *tt.initialState
				}
				return nil
			}
			tt.client.CreateFake = func(ctx context.Context, t *testing.T, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				t.Helper()
				return nil
			}
			err := operationApply(ctx, tt.object, tt.client, tt.shouldFail, nil)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
