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
	kerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_operationAssert(t *testing.T) {
	tests := []struct {
		name         string
		expected     unstructured.Unstructured
		client       *tclient.FakeClient
		expectedLogs []string
		expectErr    bool
	}{
		{
			name: "Successful match using Get",
			expected: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name": "test-pod",
					},
				},
			},
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					t.Helper()
					obj.(*unstructured.Unstructured).Object = map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "test-pod",
						},
					}
					return nil
				},
			},
			expectedLogs: []string{"ASSERT: [RUNNING...]", "ASSERT: [DONE]"},
		},
		{
			name: "Failed match using Get",
			expected: unstructured.Unstructured{
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
								"image": "test-image",
							},
						},
					},
				},
			},
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					t.Helper()
					obj.(*unstructured.Unstructured).Object = map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "Pod",
						"metadata": map[string]interface{}{
							"name": "test-pod",
						},
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "fake-container",
									"image": "fake-image",
								},
							},
						},
					}
					return nil
				},
			},
			expectErr: true,
			expectedLogs: []string{
				"ASSERT: [RUNNING...]",
				"ASSERT: [ERROR\nresource test-pod doesn't match expectation:\n    spec.containers[0].image: Invalid value: \"fake-image\": Expected value: \"test-image\"\n    spec.containers[0].name: Invalid value: \"fake-container\": Expected value: \"test-container\"]",
			},
		},
		{
			name: "Not found using Get",
			expected: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name": "test-pod",
					},
				},
			},
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					t.Helper()
					obj.(*unstructured.Unstructured).Object = nil
					return kerror.NewNotFound(schema.GroupResource{Group: "", Resource: "pods"}, "test-pod")
				},
			},
			expectErr:    true,
			expectedLogs: []string{"ASSERT: [RUNNING...]", "ASSERT: [ERROR\nactual resource not found]"},
		},
		{
			name: "Successful match using List",
			expected: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"namespace": "test-ns",
						"labels": map[string]interface{}{
							"app": "my-app",
						},
					},
				},
			},
			client: &tclient.FakeClient{
				ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					t.Helper()
					uList := list.(*unstructured.UnstructuredList)
					uList.Items = append(uList.Items, unstructured.Unstructured{
						Object: map[string]interface{}{
							"apiVersion": "apps/v1",
							"kind":       "Deployment",
							"metadata": map[string]interface{}{
								"namespace": "test-ns",
								"labels": map[string]interface{}{
									"app": "my-app",
								},
							},
						},
					})
					return nil
				},
			},
			expectedLogs: []string{"ASSERT: [RUNNING...]", "ASSERT: [DONE]"},
		},
		{
			name: "No resources found using List",
			expected: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"namespace": "test-ns",
						"labels": map[string]interface{}{
							"app": "my-app",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "test-container",
								"image": "test-image",
							},
						},
					},
				},
			},
			client: &tclient.FakeClient{
				ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					t.Helper()
					uList := list.(*unstructured.UnstructuredList)
					uList.Items = nil
					return nil
				},
			},
			expectErr:    true,
			expectedLogs: []string{"ASSERT: [RUNNING...]", "ASSERT: [ERROR\nno actual resource found]"},
		},
		{
			name: "List operation fails",
			expected: unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"namespace": "test-ns",
						"labels": map[string]interface{}{
							"app": "my-app",
						},
					},
				},
			},
			client: &tclient.FakeClient{
				ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
					t.Helper()
					return errors.New("internal server error")
				},
			},
			expectErr:    true,
			expectedLogs: []string{"ASSERT: [RUNNING...]", "ASSERT: [ERROR\ninternal server error]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			ctx := logging.IntoContext(ctxt, logger)
			assertOp := &AssertOperation{
				baseOperation: baseOperation{
					client: tt.client,
				},
				expected: tt.expected,
			}
			err := execOperation(ctx, assertOp)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
