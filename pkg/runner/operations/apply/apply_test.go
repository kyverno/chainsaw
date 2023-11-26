package apply

import (
	"context"
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_apply(t *testing.T) {
	podv1 := unstructured.Unstructured{
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
						"image": "test-image:v1",
					},
				},
			},
		},
	}
	podv2 := unstructured.Unstructured{
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
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Dry Run Resource already exists, patch it",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Resource does not exist, create it",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Dry Run Resource does not exist, create it",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		expect:      nil,
		expectedErr: nil,
	}, {
		name:   "Error while getting resource",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				return errors.New("some arbitrary error")
			},
		},
		expect:      nil,
		expectedErr: errors.New("some arbitrary error"),
	}, {
		name:   "Fail to patch existing resource",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return errors.New("patch failed")
			},
		},
		expect:      nil,
		expectedErr: errors.New("patch failed"),
	}, {
		name:   "Fail to create non-existing resource",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
				return errors.New("create failed")
			},
		},
		expect:      nil,
		expectedErr: errors.New("create failed"),
	}, {
		name:   "Unexpected patch success when should fail",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]interface{}{
					"($error != null)": true,
				},
			},
		}},
		expectedErr: errors.New("($error != null): Invalid value: false: Expected value: true"),
	}, {
		name:   "Unexpected create success when should fail",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pods").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]interface{}{
					"($error != null)": true,
				},
			},
		}},
		expectedErr: errors.New("($error != null): Invalid value: false: Expected value: true"),
	}, {
		name:   "Expected patch failure",
		object: podv2,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return errors.New("expected patch failure")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]interface{}{
					"($error)": "expected patch failure",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Expected create failure",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pods").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
				return errors.New("expected create failure")
			},
		},
		expect: []v1alpha1.Expectation{{
			Check: v1alpha1.Check{
				Value: map[string]interface{}{
					"($error)": "expected create failure",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Don't match",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: podv2.UnstructuredContent(),
			},
			Check: v1alpha1.Check{
				Value: map[string]interface{}{
					"kind": "Service",
				},
			},
		}},
		expectedErr: nil,
	}, {
		name:   "Match",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
				*obj.(*unstructured.Unstructured) = podv1
				return nil
			},
			PatchFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ctrlclient.Patch, _ ...ctrlclient.PatchOption) error {
				return nil
			},
		},
		expect: []v1alpha1.Expectation{{
			Match: &v1alpha1.Check{
				Value: podv1.UnstructuredContent(),
			},
			Check: v1alpha1.Check{
				Value: map[string]interface{}{
					"kind": "Service",
				},
			},
		}},
		expectedErr: errors.New(`kind: Invalid value: "Pod": Expected value: "Service"`),
	}, {
		name:   "Resource does not exist, create it and call cleaner",
		object: podv1,
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pods").GroupResource(), key.Name)
			},
			CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
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
			operation := operation{
				client: tt.client,
				obj:    tt.object,
				expect: tt.expect,
			}
			err := operation.Exec(ctx)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
