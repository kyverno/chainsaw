package apply

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_apply(t *testing.T) {
	podv1 := unstructured.Unstructured{
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
	podv2 := unstructured.Unstructured{
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
						"image": "test-image:v2",
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
		name:   "Resource already exists, patch it",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Dry Run Resource already exists, patch it",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Resource does not exist, create it",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Dry Run Resource does not exist, create it",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Error while getting resource",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect:      nil,
		expectedErr: errors.New("some arbitrary error"),
	}, {
		name:   "Fail to patch existing resource",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return errors.New("patch failed")
			},
		},
		expect:      nil,
		expectedErr: errors.New("patch failed"),
	}, {
		name:   "Fail to create non-existing resource",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return errors.New("create failed")
			},
		},
		expect:      nil,
		expectedErr: errors.New("create failed"),
	}, {
		name:   "Unexpected patch success when should fail",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
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
		name:   "Unexpected create success when should fail",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pods").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
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
		name:   "Expected patch failure",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return errors.New("expected patch failure")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]any{
					"($error)": "expected patch failure",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Expected create failure",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pods").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return errors.New("expected create failure")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]any{
					"($error)": "expected create failure",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Don't match",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: podv2.UnstructuredContent(),
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
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ client.Object, _ client.Patch, _ ...client.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: podv1.UnstructuredContent(),
			},
			Check: v1alpha1.Check{
				Value: map[string]any{
					"kind": "Service",
				},
			},
		}},
		expectedErr: errors.New(`kind: Invalid value: "Pod": Expected value: "Service"`),
	}, {
		name:   "Resource does not exist, create it and call cleaner",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pods").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctx := logging.IntoContext(context.TODO(), logger)
			toCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			ctx = toCtx
			operation := New(
				tt.client,
				tt.object,
				nil,
				nil,
				false,
				tt.expect,
				nil,
			)
			outputs, err := operation.Exec(ctx, nil)
			assert.Nil(t, outputs)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
