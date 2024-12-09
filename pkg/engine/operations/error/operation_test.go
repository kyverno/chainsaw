package error

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/client"
	tclient "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	tlogging "github.com/kyverno/chainsaw/pkg/engine/logging/testing"
	"github.com/kyverno/chainsaw/pkg/engine/namespacer"
	tnamespacer "github.com/kyverno/chainsaw/pkg/engine/namespacer/testing"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Test_operationError(t *testing.T) {
	expected := unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
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
			ListFn: func(ctx context.Context, _ int, list client.ObjectList, opts ...client.ListOption) error {
				return kerrors.NewNotFound(list.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
			},
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerrors.NewNotFound(obj.GetObjectKind().GroupVersionKind().GroupVersion().WithResource("pod").GroupResource(), "test-pod")
			},
		},
		expectedErr:  nil,
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: DONE - []"},
	}, {
		name:     "Internal error",
		expected: expected,
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("internal error")
			},
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return errors.New("internal error")
			},
		},
		expectedErr:  errors.New("internal error"),
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: ERROR - [=== ERROR\ninternal error]"},
	}, {
		name:     "Resource matches actual",
		expected: expected,
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list client.ObjectList, opts ...client.ListOption) error {
				return nil
			},
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
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
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name": "test-pod",
				},
				"spec": map[string]any{
					"(foo('bar'))": "test-pod",
				},
			},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				t.Helper()
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod",
					},
					"spec": map[string]any{},
				}
				return nil
			},
		},
		expectedErr:  fmt.Errorf("spec.(foo('bar')): Internal error: unknown function: foo"),
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: ERROR - [=== ERROR\nspec.(foo('bar')): Internal error: unknown function: foo]"},
	}, {
		name: "No resources found using List",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"namespace": "test-ns",
					"labels": map[string]any{
						"app": "my-app",
					},
				},
				"spec": map[string]any{
					"containers": []any{
						map[string]any{
							"name":  "test-container",
							"image": "test-image",
						},
					},
				},
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list client.ObjectList, opts ...client.ListOption) error {
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
			Object: map[string]any{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]any{
					"labels": map[string]any{
						"app": "my-app",
					},
				},
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list client.ObjectList, opts ...client.ListOption) error {
				assert.Contains(t, opts, client.InNamespace("bar"))
				uList := list.(*unstructured.UnstructuredList)
				uList.Items = nil
				return nil
			},
			IsObjectNamespacedFn: func(int, runtime.Object) (bool, error) {
				return true, nil
			},
		},
		namespacer: func(c client.Client) namespacer.Namespacer {
			return namespacer.New("bar")
		},
		expectedLogs: []string{"ERROR: RUN - []", "ERROR: DONE - []"},
	}, {
		name: "with namespacer error",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]any{
					"labels": map[string]any{
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
				ApplyFn: func(call int, client client.Client, obj client.Object) error {
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
				apis.DefaultCompilers,
				tt.client,
				tt.expected,
				nspacer,
				false,
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
