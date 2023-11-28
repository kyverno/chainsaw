package error

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
	"github.com/kyverno/chainsaw/pkg/runner/namespacer"
	tnamespacer "github.com/kyverno/chainsaw/pkg/runner/namespacer/testing"
	ttesting "github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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
		namespacer   func(c client.Client) namespacer.Namespacer
		expectedErr  error
		expectedLogs []string
	}{{
		name:     "Resource not found",
		expected: expected,
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				return kerrors.NewNotFound(list.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
			},
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: DONE - []"},
	}, {
		name:     "Internal error",
		expected: expected,
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				return errors.New("internal error")
			},
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return errors.New("internal error")
			},
		},
		expectedErr:  errors.New("internal error"),
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: ERROR - [=== ERROR\ninternal error]"},
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
		expectedErr:  fmt.Errorf("v1/Pod/foo/test-pod - resource matches expectation"),
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: ERROR - [=== ERROR\nv1/Pod/foo/test-pod - resource matches expectation]"},
	}, {
		name: "Bad assert",
		expected: unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name": "test-pod",
				},
				"spec": map[string]interface{}{
					"(foo('bar'))": "test-pod",
				},
			},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
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
		expectedErr:  fmt.Errorf("spec.(foo('bar')): Internal error: unknown function: foo"),
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: ERROR - [=== ERROR\nspec.(foo('bar')): Internal error: unknown function: foo]"},
	}, {
		name: "No resources found using List",
		expected: unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"namespace": "test-ns",
					"labels": map[string]interface{}{
						"app": "my-app",
					},
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "test-container",
							"image": "test-image",
						},
					},
				},
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				t.Helper()
				uList := list.(*unstructured.UnstructuredList)
				uList.Items = nil
				return nil
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: DONE - []"},
	}, {
		name: "with namespacer",
		expected: unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": "my-app",
					},
				},
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				t := ttesting.FromContext(ctx)
				assert.Contains(t, opts, ctrlclient.InNamespace("bar"))
				uList := list.(*unstructured.UnstructuredList)
				uList.Items = nil
				return nil
			},
			IsObjectNamespacedFn: func(int, runtime.Object) (bool, error) {
				return true, nil
			},
		},
		namespacer: func(c client.Client) namespacer.Namespacer {
			return namespacer.New(c, "bar")
		},
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: DONE - []"},
	}, {
		name: "with namespacer error",
		expected: unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": "my-app",
					},
				},
			},
		},
		client: &tclient.FakeClient{
			IsObjectNamespacedFn: func(int, runtime.Object) (bool, error) {
				return true, nil
			},
		},
		namespacer: func(c client.Client) namespacer.Namespacer {
			return &tnamespacer.FakeNamespacer{
				ApplyFn: func(obj ctrlclient.Object, call int) error {
					return errors.New("namespacer error")
				},
			}
		},
		expectedErr:  errors.New("namespacer error"),
		expectedLogs: []string{"ERROR: ERROR - [=== ERROR\nnamespacer error]"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			var nspacer namespacer.Namespacer
			if tt.namespacer != nil {
				nspacer = tt.namespacer(tt.client)
			}
			operation := New(
				tt.client,
				tt.expected,
				nspacer,
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
