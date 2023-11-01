package operations

import (
	"context"
	"testing"

	fakeClient "github.com/kyverno/chainsaw/pkg/runner/client"
	fakeLogger "github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAssert(t *testing.T) {
	tests := []struct {
		name         string
		expected     unstructured.Unstructured
		fakeClient   *fakeClient.FakeClient
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
			fakeClient: &fakeClient.FakeClient{
				T: t,
				GetFake: func(ctx context.Context, t *testing.T, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &fakeLogger.MockLogger{}
			err := Assert(context.TODO(), mockLogger, tt.expected, tt.fakeClient)

			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.expectedLogs, mockLogger.Logs)
		})
	}
}
