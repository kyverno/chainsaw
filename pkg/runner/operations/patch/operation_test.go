package patch

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_create(t *testing.T) {
	pod := unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name": "test-pod",
			},
			"spec": map[string]any{
				"containers": []any{
					map[string]any{
						"name":  "test-container",
						"image": "test-image:v1",
					},
				},
			},
		},
	}
	tests := []struct {
		name        string
		object      unstructured.Unstructured
		client      *tclient.FakeClient
		expect      []v1alpha1.Expectation
		expectedErr error
	}{{
		name:   "Resource doesn't exist",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
		},
		expect:      nil,
		expectedErr: errors.New(`pod "test-pod" not found`),
	}, {
		name:   "Dry Run Resource doesn't exist",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
		},
		expect:      nil,
		expectedErr: errors.New(`pod "test-pod" not found`),
	}, {
		name:   "Resource exists, patch it",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Dry Run Resource exists, patch it",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "failed get",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, _ ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect:      nil,
		expectedErr: errors.New("some arbitrary error"),
	}, {
		name:   "failed patch",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect:      nil,
		expectedErr: errors.New("some arbitrary error"),
	}, {
		name:   "failed patch (expected)",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]any{
					"($error)": "some arbitrary error",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Should fail is true but no error occurs",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]any{
					"($error != null)": true,
				},
			},
		}},
		expectedErr: errors.New("($error != null): Invalid value: false: Expected value: true"),
	}, {
		name:   "Don't match",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: map[string]any{
					"foo": "bar",
				},
			},
			Check: v1alpha1.Check{
				Value: map[string]any{
					"kind": "Service",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Match",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: pod.UnstructuredContent(),
			},
			Check: v1alpha1.Check{
				Value: map[string]any{
					"kind": "Service",
				},
			},
		}},
		expectedErr: errors.New(`kind: Invalid value: "Pod": Expected value: "Service"`),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctx := logging.IntoContext(context.TODO(), logger)
			toCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			ctx = toCtx
			operation := New(
				tt.client,
				tt.object,
				nil,
				nil,
				false,
				tt.expect,
			)
			err := operation.Exec(ctx)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
