package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/cleanup/cleaner"
	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/mocks"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
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
		cleaner     cleaner.Cleaner
		expect      []v1alpha1.Expectation
		expectedErr error
	}{{
		name:   "Resource already exists",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return kerrors.NewAlreadyExists(
					schema.GroupResource{Group: "", Resource: "pods"},
					"test-pod",
				)
			},
		},
		expect:      nil,
		expectedErr: errors.New("the resource already exists in the cluster"),
	}, {
		name:   "Dry Run Resource already exists",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
				*obj.(*unstructured.Unstructured) = pod
				return nil
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return kerrors.NewAlreadyExists(
					schema.GroupResource{Group: "", Resource: "pods"},
					"test-pod",
				)
			},
		},
		expect:      nil,
		expectedErr: errors.New("the resource already exists in the cluster"),
	}, {
		name:   "Resource does not exist, create it",
		object: pod,
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
		object: pod,
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
		name:   "failed get",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ client.ObjectKey, _ client.Object, opts ...client.GetOption) error {
				return errors.New("some arbitrary error")
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return errors.New("unexpected create call")
			},
		},
		expect:      nil,
		expectedErr: errors.New("some arbitrary error"),
	}, {
		name:   "failed create",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect:      nil,
		expectedErr: errors.New("some arbitrary error"),
	}, {
		name:   "failed create (expected)",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.NewCheck(
				map[string]any{
					"($error)": "some arbitrary error",
				},
			),
		}},
		expectedErr: nil,
	}, {
		name:   "Cleaner function executed",
		object: pod,
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
		name:   "Should fail is true but no error occurs",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.NewCheck(
				map[string]any{
					"($error != null)": true,
				},
			),
		}},
		expectedErr: errors.New("($error != null): Invalid value: false: Expected value: true"),
	}, {
		name:   "Don't match",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: ptr.To(v1alpha1.NewMatch(
				map[string]any{
					"foo": "bar",
				},
			)),
			Check: v1alpha1.NewCheck(
				map[string]any{
					"kind": "Service",
				},
			),
		}},
		expectedErr: nil,
	}, {
		name:   "Match",
		object: pod,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ client.Object, _ ...client.CreateOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: ptr.To(v1alpha1.NewMatch(pod.UnstructuredContent())),
			Check: v1alpha1.NewCheck(
				map[string]any{
					"kind": "Service",
				},
			),
		}},
		expectedErr: errors.New(`kind: Invalid value: "Pod": Expected value: "Service"`),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.Logger{}
			ctx := logging.WithLogger(context.TODO(), logger)
			toCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			ctx = toCtx
			operation := New(
				apis.DefaultCompilers,
				tt.client,
				tt.object,
				nil,
				tt.cleaner,
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

func Test_retry_logic(t *testing.T) {
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
		cleaner     cleaner.Cleaner
		expect      []v1alpha1.Expectation
		expectedErr error
	}{
		{
			name:   "conflict error should be retried and eventually succeed",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, call int, _ client.Object, _ ...client.CreateOption) error {
					if call < 2 {
						return kerrors.NewConflict(
							schema.GroupResource{Group: "", Resource: "pods"},
							"test-pod",
							errors.New("conflict error"),
						)
					}
					return nil
				},
			},
			expect:      nil,
			expectedErr: nil,
		},
		{
			name:   "server timeout error should be retried and eventually succeed",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, call int, _ client.Object, _ ...client.CreateOption) error {
					if call < 2 {
						return kerrors.NewServerTimeout(
							schema.GroupResource{Group: "", Resource: "pods"},
							"create",
							10,
						)
					}
					return nil
				},
			},
			expect:      nil,
			expectedErr: nil,
		},
		{
			name:   "too many requests error should be retried and eventually succeed",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, call int, _ client.Object, _ ...client.CreateOption) error {
					if call < 2 {
						return kerrors.NewTooManyRequests(
							"too many requests",
							10,
						)
					}
					return nil
				},
			},
			expect:      nil,
			expectedErr: nil,
		},
		{
			name:   "service unavailable error should be retried and eventually succeed",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, call int, _ client.Object, _ ...client.CreateOption) error {
					if call < 2 {
						return kerrors.NewServiceUnavailable("service unavailable")
					}
					return nil
				},
			},
			expect:      nil,
			expectedErr: nil,
		},
		{
			name:   "already exists error should not be retried",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, call int, _ client.Object, _ ...client.CreateOption) error {
					return kerrors.NewAlreadyExists(
						schema.GroupResource{Group: "", Resource: "pods"},
						"test-pod",
					)
				},
			},
			expect:      nil,
			expectedErr: errors.New("the resource already exists in the cluster"),
		},
		{
			name:   "permanent error should not be retried",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, call int, _ client.Object, _ ...client.CreateOption) error {
					return kerrors.NewBadRequest("bad request error")
				},
			},
			expect:      nil,
			expectedErr: errors.New("bad request error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &mocks.Logger{}
			ctx := logging.WithLogger(context.TODO(), logger)
			toCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			ctx = toCtx
			operation := New(
				apis.DefaultCompilers,
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
