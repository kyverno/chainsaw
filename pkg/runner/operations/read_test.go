package operations

// import (
// 	"context"
// 	"testing"

// 	fakeClient "github.com/kyverno/chainsaw/pkg/runner/client"
// 	"github.com/stretchr/testify/assert"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
// )

// func TestRead(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		expected       ctrlclient.Object
// 		client         *fakeClient.FakeClient
// 		expectedResult []unstructured.Unstructured
// 		expectedError  string
// 	}{
// 		{
// 			name: "Test Get",
// 			expected: &unstructured.Unstructured{
// 				Object: map[string]interface{}{
// 					"apiVersion": "v1",
// 					"kind":       "Pod",
// 					"metadata": map[string]interface{}{
// 						"name": "test-pod",
// 					},
// 				},
// 			},
// 			client: &fakeClient.FakeClient{
// 				T: t,
// 				GetFake: func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 					t.Helper()
// 					obj.(*unstructured.Unstructured).Object = map[string]interface{}{
// 						"apiVersion": "v1",
// 						"kind":       "Pod",
// 						"metadata": map[string]interface{}{
// 							"name": "test-pod",
// 						},
// 					}
// 					return nil
// 				},
// 			},
// 			expectedResult: []unstructured.Unstructured{
// 				{
// 					Object: map[string]interface{}{
// 						"apiVersion": "v1",
// 						"kind":       "Pod",
// 						"metadata": map[string]interface{}{
// 							"name": "test-pod",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "Test List",
// 			expected: &unstructured.Unstructured{
// 				Object: map[string]interface{}{
// 					"apiVersion": "v1",
// 					"kind":       "Pod",
// 				},
// 			},
// 			client: &fakeClient.FakeClient{
// 				T: t,
// 				ListFake: func(ctx context.Context, t *testing.T, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 					t.Helper()
// 					list.(*unstructured.UnstructuredList).Items = []unstructured.Unstructured{
// 						{
// 							Object: map[string]interface{}{
// 								"apiVersion": "v1",
// 								"kind":       "Pod",
// 								"metadata": map[string]interface{}{
// 									"name": "test-pod-1",
// 								},
// 							},
// 						},
// 						{
// 							Object: map[string]interface{}{
// 								"apiVersion": "v1",
// 								"kind":       "Pod",
// 								"metadata": map[string]interface{}{
// 									"name": "test-pod-2",
// 								},
// 							},
// 						},
// 					}
// 					return nil
// 				},
// 			},
// 			expectedResult: []unstructured.Unstructured{
// 				{
// 					Object: map[string]interface{}{
// 						"apiVersion": "v1",
// 						"kind":       "Pod",
// 						"metadata": map[string]interface{}{
// 							"name": "test-pod-1",
// 						},
// 					},
// 				},
// 				{
// 					Object: map[string]interface{}{
// 						"apiVersion": "v1",
// 						"kind":       "Pod",
// 						"metadata": map[string]interface{}{
// 							"name": "test-pod-2",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := read(context.TODO(), tt.expected, tt.client)
// 			if tt.expectedError == "" {
// 				assert.Nil(t, err)
// 			} else {
// 				assert.EqualError(t, err, tt.expectedError)
// 			}
// 			assert.Equal(t, tt.expectedResult, result)
// 		})
// 	}
// }
