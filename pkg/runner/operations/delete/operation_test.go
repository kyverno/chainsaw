package delete

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	ttesting "github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
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
			GetFn: func(_ context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}, {
		name:   "failed get",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(_ context.Context, _ int, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return kerrors.NewInternalError(errors.New("failed to get the pod"))
			},
			DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  kerrors.NewInternalError(errors.New("failed to get the pod")),
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: ERROR - [=== ERROR\nInternal error occurred: failed to get the pod]"},
	}, {
		name:   "failed delete",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(_ context.Context, _ int, _ ctrlclient.ObjectKey, _ ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return nil
			},
			DeleteFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
				return kerrors.NewInternalError(errors.New("failed to delete the pod"))
			},
		},
		expectedErr:  kerrors.NewInternalError(errors.New("failed to delete the pod")),
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: ERROR - [=== ERROR\nInternal error occurred: failed to delete the pod]"},
	}, {
		name:   "ok",
		object: pod,
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
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}, {
		name:   "poll succeeds but returns error after",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return nil
			},
			DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
				return nil
			},
		},
		expectedErr:  context.DeadlineExceeded,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: ERROR - [=== ERROR\ncontext deadline exceeded]"},
	}, {
		name:   "with namespacer",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				t := ttesting.FromContext(ctx)
				assert.Equal(t, "bar", key.Namespace)
				if call == 0 {
					return nil
				}
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
				t := ttesting.FromContext(ctx)
				assert.Equal(t, "bar", obj.GetNamespace())
				return nil
			},
			IsObjectNamespacedFn: func(int, runtime.Object) (bool, error) {
				return true, nil
			},
		},
		namespacer: func(c client.Client) namespacer.Namespacer {
			return namespacer.New(c, "bar")
		},
		expectedErr:  nil,
		expectedLogs: []string{"DELETE: RUN - []", "DELETE: DONE - []"},
	}, {
		name:   "with check",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				if call < 10 {
					return nil
				}
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, _ ...ctrlclient.DeleteOption) error {
				t := ttesting.FromContext(ctx)
				assert.Equal(t, 1, call)
				return errors.New("dummy error")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]any{
					"($error == 'dummy error')": true,
				},
			},
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
				tt.client,
				tt.object,
				nspacer,
				nil,
				false,
				tt.expect...,
			)
			logger := &tlogging.FakeLogger{}
			err := operation.Exec(ttesting.IntoContext(logging.IntoContext(ctx, logger), t))
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
