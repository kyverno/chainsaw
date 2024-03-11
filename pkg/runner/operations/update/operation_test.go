package update

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/cleanup"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)
 
func Test_update(t *testing.T) {
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
        name:   "Resource exists, update it",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
                *obj.(*unstructured.Unstructured) = pod
                return nil
            },
            UpdateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.UpdateOption) error {
                return nil
            },
        },
        expect:      nil,
        expectedErr: nil,
    },{
        name:   "Error while updating resource",
        object: podv1,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, _ ctrlclient.ObjectKey, obj ctrlclient.Object, _ ...ctrlclient.GetOption) error {
                *obj.(*unstructured.Unstructured) = pod
                return nil
            },
            UpdateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.UpdateOption) error {
                return errors.New("update failed")
            },
        },
        expect:      nil,
        expectedErr: errors.New("update failed"),
    },{
        name:   "Dry Run Resource exists, update it",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
                return nil
            },
        },
        expectedErr: nil,
        dryRun:      true, 
    }, {
        name:   "Resource does not exist, update it",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
                return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
            },
            UpdateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.UpdateOption) error {
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
        name:   "failed update (expected)",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
                return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
            },
            UpdateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.UpdateOption) error {
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
    },
    {
        name:   "Cleanup function executed on update failure",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
                return nil
            },
            UpdateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.UpdateOption) error {
                return errors.New("update failed")
            },
        },
        expectedErr:   errors.New("update failed"),
    }, {
        name:   "Should fail is true but no error occurs",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
                return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), key.Name)
            },
            UpdateFn: func(_ context.Context, _ int, _ ctrlclient.Object, _ ...ctrlclient.UpdateOption) error {
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
                return nil
            },
        },
        expect: []v1alpha1.Expectation{
            {Key: "metadata.name", Operator: "Equals", Value: "unexpected-pod"},
        },
        expectedErr: nil,
    }, {
        name:   "Match",
        object: pod,
        client: &tclient.FakeClient{
            GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
               return nil
            },
        },
        expect: []v1alpha1.Expectation{
            {Key: "metadata.name", Operator: "Equals", Value: "test-pod"},
        },
        expectedErr: errors.New("an invalid value was found"),
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
