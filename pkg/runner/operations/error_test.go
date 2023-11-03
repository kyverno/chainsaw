package operations

// import (
// 	"context"
// 	"testing"

// 	fakeClient "github.com/kyverno/chainsaw/pkg/runner/client"
// 	"github.com/kyverno/chainsaw/pkg/runner/logging"
// 	"github.com/stretchr/testify/assert"
// 	"k8s.io/apimachinery/pkg/api/errors"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
// )

// func Test_operationError(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		expected     unstructured.Unstructured
// 		client       *fakeClient.FakeClient
// 		expectedErr  error
// 		expectedLogs []string
// 	}{{
// 		name: "Resource not found",
// 		expected: unstructured.Unstructured{
// 			Object: map[string]interface{}{
// 				"apiVersion": "v1",
// 				"kind":       "Pod",
// 				"metadata": map[string]interface{}{
// 					"name": "test-pod",
// 				},
// 			},
// 		},
// 		client: &fakeClient.FakeClient{
// 			T: &testing.T{},
// 			ListFake: func(ctx context.Context, t *testing.T, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 				t.Helper()
// 				return errors.NewNotFound(list.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
// 			},
// 			GetFake: func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				t.Helper()
// 				return errors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
// 			},
// 		},
// 		expectedErr:  nil,
// 		expectedLogs: []string{"ERROR : [RUNNING...]", "ERROR : [DONE]"},
// 		// }, {
// 		// 	name: "Resource matches actual",
// 		// 	expected: unstructured.Unstructured{
// 		// 		Object: map[string]interface{}{
// 		// 			"apiVersion": "v1",
// 		// 			"kind":       "Pod",
// 		// 			"metadata": map[string]interface{}{
// 		// 				"name": "test-pod",
// 		// 			},
// 		// 		},
// 		// 	},
// 		// 	client: &fakeClient.FakeClient{
// 		// 		T: &testing.T{},
// 		// 		ListFake: func(ctx context.Context, t *testing.T, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 		// 			return nil
// 		// 		},
// 		// 		GetFake: func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 		// 			t.Helper()
// 		// 			uObj, ok := obj.(*unstructured.Unstructured)
// 		// 			if !ok {
// 		// 				t.Fatalf("obj is not of type *unstructured.Unstructured, it's %T", obj)
// 		// 			}
// 		// 			uObj.Object = map[string]interface{}{
// 		// 				"apiVersion": "v1",
// 		// 				"kind":       "Pod",
// 		// 				"metadata": map[string]interface{}{
// 		// 					"name": "test-pod",
// 		// 				},
// 		// 			}
// 		// 			return nil
// 		// 		},
// 		// 	},
// 		// 	expectedErr:  fmt.Errorf("found an actual resource matching expectation (v1/Pod / test-pod)"),
// 		// 	expectedLogs: []string{"ERROR : [RUNNING...]", "ERROR : [ERROR\nfound an actual resource matching expectation (v1/Pod / test-pod)]"},
// 	}}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := &MockLogger{}
// 			ctx := logging.IntoContext(context.TODO(), logger)
// 			err := operationError(ctx, tt.expected, tt.client)
// 			if tt.expectedErr != nil {
// 				assert.EqualError(t, err, tt.expectedErr.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.Equal(t, tt.expectedLogs, logger.Logs)
// 		})
// 	}
// }
