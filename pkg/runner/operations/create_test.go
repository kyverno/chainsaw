package operations

import (
	"context"
	"errors"
	"testing"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_create(t *testing.T) {
	var cleanerCalled bool
	testCleaner := func(obj ctrlclient.Object, c client.Client) {
		cleanerCalled = true
	}
	pod := &unstructured.Unstructured{
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
	tests := []struct {
		name        string
		object      ctrlclient.Object
		client      *tclient.FakeClient
		cleaner     cleanup.Cleaner
		shouldFail  bool
		dryrun      bool
		expectedErr error
	}{
		{
			name:   "Resource already exists",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					*obj.(*unstructured.Unstructured) = *pod.DeepCopy()
					return nil
				},
			},
			shouldFail:  false,
			expectedErr: errors.New("the resource already exists in the cluster"),
		},
		{
			name:   "Dry Run Resource already exists",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
					*obj.(*unstructured.Unstructured) = *pod.DeepCopy()
					return nil
				},
			},
			shouldFail:  false,
			expectedErr: errors.New("the resource already exists in the cluster"),
			dryrun:      true,
		},
		{
			name:   "Resource does not exist, create it",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					return nil
				},
			},
			shouldFail:  false,
			expectedErr: nil,
		},
		{
			name:   "Dry Run Resource does not exist, create it",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					return nil
				},
			},
			shouldFail:  false,
			expectedErr: nil,
			dryrun:      true,
		},
		{
			name:   "failed get",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, _ ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return errors.New("some arbitrary error")
				},
			},
			shouldFail:  false,
			expectedErr: errors.New("some arbitrary error"),
		},
		{
			name:   "failed create",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					return errors.New("some arbitrary error")
				},
			},
			shouldFail:  false,
			expectedErr: errors.New("some arbitrary error"),
		},
		{
			name:   "failed create (expected)",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					return errors.New("some arbitrary error")
				},
			},
			shouldFail:  true,
			expectedErr: nil,
		},
		{
			name:   "Cleaner function executed",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					return nil
				},
			},
			shouldFail:  false,
			expectedErr: nil,
			cleaner:     testCleaner,
		},
		{
			name:   "Should fail is true but no error occurs",
			object: pod,
			client: &tclient.FakeClient{
				GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
				},
				CreateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.CreateOption) error {
					return nil
				},
			},
			shouldFail:  true,
			expectedErr: errors.New("an error was expected but didn't happen"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &tlogging.FakeLogger{}
			ctx := logging.IntoContext(context.TODO(), logger)
			createOp := &CreateOperation{
				baseOperation: baseOperation{
					client: tt.client,
				},
				obj:        tt.object,
				dryRun:     tt.dryrun,
				cleaner:    tt.cleaner,
				shouldFail: tt.shouldFail,
			}
			err := createOp.Exec(ctx)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.cleaner != nil {
				assert.True(t, cleanerCalled, "cleaner was not called when expected")
			}
		})
	}
}
