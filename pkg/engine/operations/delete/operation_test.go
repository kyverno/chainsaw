package delete

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_operationDelete(t *testing.T) {
	pod := unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name": "test-pod",
			},
		},
	}
	tests := []struct {
		name         string
		object       unstructured.Unstructured
		client       *tclient.FakeClient
		namespacer   func(c client.Client) namespacer.Namespacer
		expect       []v1alpha1.Expectation
		expectedErr  error
		expectedLogs []string
	}{{
		name:   "not found",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(_ context.Context, _ int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(_ context.Context, _ int, _ client.Object, _ ...client.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}, {
		name:   "failed get",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(_ context.Context, _ int, _ client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
				return kerrors.NewInternalError(errors.New("failed to get the pod"))
			},
			DeleteFn: func(_ context.Context, _ int, _ client.Object, _ ...client.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  kerrors.NewInternalError(errors.New("failed to get the pod")),
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: ERROR - [=== ERROR\nInternal error occurred: failed to get the pod]"},
	}, {
		name:   "failed delete",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(_ context.Context, _ int, _ client.ObjectKey, _ client.Object, _ ...client.GetOption) error {
				return nil
			},
			DeleteFn: func(_ context.Context, _ int, _ client.Object, _ ...client.DeleteOption) error {
				return kerrors.NewInternalError(errors.New("failed to delete the pod"))
			},
		},
		expectedErr:  kerrors.NewInternalError(errors.New("failed to delete the pod")),
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: ERROR - [=== ERROR\nInternal error occurred: failed to delete the pod]"},
	}, {
		name:   "ok",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(_ context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				if call < 10 {
					return nil
				}
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(_ context.Context, _ int, _ client.Object, _ ...client.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}, {
		name:   "poll succeeds but returns error after",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				return nil
			},
			DeleteFn: func(ctx context.Context, call int, obj client.Object, _ ...client.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  context.DeadlineExceeded,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: ERROR - [=== ERROR\ncontext deadline exceeded]"},
	}, {
		name:   "with namespacer",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				assert.Equal(t, "bar", key.Namespace)
				if call == 0 {
					return nil
				}
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(ctx context.Context, call int, obj client.Object, _ ...client.DeleteOption) error {
				assert.Equal(t, "bar", obj.GetNamespace())
				return nil
			},
			IsObjectNamespacedFn: func(int, runtime.Object) (bool, error) {
				return true, nil
			},
		},
		namespacer: func(c client.Client) namespacer.Namespacer {
			return namespacer.New("bar")
		},
		expectedErr:  nil,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}, {
		name:   "with check",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				if call < 10 {
					return nil
				}
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(ctx context.Context, call int, obj client.Object, _ ...client.DeleteOption) error {
				assert.Equal(t, 1, call)
				return errors.New("dummy error")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.NewCheck(
				map[string]any{
					"($error == 'dummy error')": true,
				},
			),
		}},
		expectedErr:  nil,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			var nspacer namespacer.Namespacer
			if tt.namespacer != nil {
				nspacer = tt.namespacer(tt.client)
			}
			operation := New(
				apis.DefaultCompilers,
				tt.client,
				tt.object,
				nspacer,
				false,
				metav1.DeletePropagationForeground,
				tt.expect...,
			)
			logger := &tlogging.FakeLogger{}
			outputs, err := operation.Exec(logging.IntoContext(ctx, logger), nil)
			assert.Nil(t, outputs)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
