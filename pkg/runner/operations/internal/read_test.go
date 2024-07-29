package internal

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name           string
		expected       client.Object
		client         *tclient.FakeClient
		expectedResult []unstructured.Unstructured
		expectedError  string
	}{{
		name: "Test Get",
		expected: &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name": "test-pod",
				},
			},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod",
					},
				}
				return nil
			},
		},
		expectedResult: []unstructured.Unstructured{
			{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod",
					},
				},
			},
		},
	}, {
		name: "Test List",
		expected: &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list client.ObjectList, opts ...client.ListOption) error {
				list.(*unstructured.UnstructuredList).Items = []unstructured.Unstructured{
					{
						Object: map[string]any{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]any{
								"name": "test-pod-1",
							},
						},
					},
					{
						Object: map[string]any{
							"apiVersion": "v1",
							"kind":       "Pod",
							"metadata": map[string]any{
								"name": "test-pod-2",
							},
						},
					},
				}
				return nil
			},
		},
		expectedResult: []unstructured.Unstructured{
			{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod-1",
					},
				},
			},
			{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod-2",
					},
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Read(context.TODO(), tt.expected, tt.client)
			if tt.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
