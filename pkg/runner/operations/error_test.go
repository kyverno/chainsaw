package operations

import (
	"context"
	"fmt"
	"testing"
	"time"

	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_operationError(t *testing.T) {
	expected := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"namespace": "foo",
				"name":      "test-pod",
			},
		},
	}
	tests := []struct {
		name         string
		expected     unstructured.Unstructured
		client       *tclient.FakeClient
		expectedErr  error
		expectedLogs []string
	}{{
		name:     "Resource not found",
		expected: expected,
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				return errors.NewNotFound(list.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
			},
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return errors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"ERROR : [RUNNING...]", "ERROR : [DONE]"},
	}, {
		name:     "Resource matches actual",
		expected: expected,
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				return nil
			},
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				uObj, ok := obj.(*unstructured.Unstructured)
				if !ok {
					t.Fatalf("obj is not of type *unstructured.Unstructured, it's %T", obj)
				}
				uObj.Object = expected.Object
				return nil
			},
		},
		expectedErr:  fmt.Errorf("found an actual resource matching expectation (v1/Pod / foo/test-pod)"),
		expectedLogs: []string{"ERROR : [RUNNING...]", "ERROR : [ERROR\nfound an actual resource matching expectation (v1/Pod / foo/test-pod)]"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctx, cancel := context.WithTimeout(logging.IntoContext(context.TODO(), logger), 1*time.Second)
			defer cancel()
			err := operationError(ctx, tt.expected, tt.client)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
