package assert

import (
	"context"
	"errors"
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
	kerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Test_operationAssert(t *testing.T) {
	tests := []struct {
		name         string
		expected     unstructured.Unstructured
		client       *tclient.FakeClient
		namespacer   func(c client.Client) namespacer.Namespacer
		expectedLogs []string
		expectErr    bool
	}{{
		name: "Successful match using Get",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name": "test-pod",
				},
			},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				t.Helper()
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod",
					},
				}
				return nil
			},
		},
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: DONE - []"},
	}, {
		name: "Failed match using Get",
		expected: unstructured.Unstructured{
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
							"image": "test-image",
						},
					},
				},
			},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				t.Helper()
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod",
					},
					"spec": map[string]any{
						"containers": []any{
							map[string]any{
								"name":  "fake-container",
								"image": "fake-image",
							},
						},
					},
				}
				return nil
			},
		},
		expectErr: true,
		expectedLogs: []string{
			"ASSERT: RUN - []",
			"ASSERT: ERROR - [=== ERROR\n---------------\nv1/Pod/test-pod\n---------------\n* spec.containers[0].image: Invalid value: \"fake-image\": Expected value: \"test-image\"\n* spec.containers[0].name: Invalid value: \"fake-container\": Expected value: \"test-container\"\n\n--- expected\n+++ actual\n@@ -4,6 +4,6 @@\n   name: test-pod\n spec:\n   containers:\n-  - image: test-image\n-    name: test-container\n+  - image: fake-image\n+    name: fake-container]",
		},
	}, {
		name: "Not found using Get",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name": "test-pod",
				},
			},
		},
		client: &tclient.FakeClient{
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				t.Helper()
				obj.(*unstructured.Unstructured).Object = nil
				return kerror.NewNotFound(schema.GroupResource{Group: "", Resource: "pods"}, "test-pod")
			},
		},
		expectErr:    true,
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: ERROR - [=== ERROR\nactual resource not found]"},
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
			GetFn: func(ctx context.Context, _ int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				t.Helper()
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "test-pod",
					},
				}
				return nil
			},
		},
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: ERROR - [=== ERROR\nspec.(foo('bar')): Internal error: unknown function: foo]"},
		expectErr:    true,
	}, {
		name: "Successful match using List",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]any{
					"namespace": "test-ns",
					"labels": map[string]any{
						"app": "my-app",
					},
				},
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				t.Helper()
				uList := list.(*unstructured.UnstructuredList)
				uList.Items = append(uList.Items, unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]any{
							"namespace": "test-ns",
							"labels": map[string]any{
								"app": "my-app",
							},
						},
					},
				})
				return nil
			},
		},
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: DONE - []"},
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
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				t.Helper()
				uList := list.(*unstructured.UnstructuredList)
				uList.Items = nil
				return nil
			},
		},
		expectErr:    true,
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: ERROR - [=== ERROR\nno actual resource found]"},
	}, {
		name: "List operation fails",
		expected: unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]any{
					"namespace": "test-ns",
					"labels": map[string]any{
						"app": "my-app",
					},
				},
			},
		},
		client: &tclient.FakeClient{
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				t.Helper()
				return errors.New("internal server error")
			},
		},
		expectErr:    true,
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: ERROR - [=== ERROR\ninternal server error]"},
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
			ListFn: func(ctx context.Context, _ int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
				t := ttesting.FromContext(ctx)
				assert.Contains(t, opts, ctrlclient.InNamespace("bar"))
				uList := list.(*unstructured.UnstructuredList)
				uList.Items = append(uList.Items, unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "apps/v1",
						"kind":       "Deployment",
						"metadata": map[string]any{
							"namespace": "bar",
							"labels": map[string]any{
								"app": "my-app",
							},
						},
					},
				})
				return nil
			},
			IsObjectNamespacedFn: func(int, runtime.Object) (bool, error) {
				return true, nil
			},
		},
		namespacer: func(c client.Client) namespacer.Namespacer {
			return namespacer.New(c, "bar")
		},
		expectedLogs: []string{"ASSERT: RUN - []", "ASSERT: DONE - []"},
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
				ApplyFn: func(obj ctrlclient.Object, call int) error {
					return errors.New("namespacer error")
				},
			}
		},
		expectErr:    true,
		expectedLogs: []string{"ASSERT: ERROR - [=== ERROR\nnamespacer error]"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var nspacer namespacer.Namespacer
			if tt.namespacer != nil {
				nspacer = tt.namespacer(tt.client)
			}
			operation := New(
				tt.client,
				tt.expected,
				nspacer,
				nil,
			)
			logger := &tlogging.FakeLogger{}
			err := operation.Exec(ttesting.IntoContext(logging.IntoContext(ctx, logger), t))
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expectedLogs, logger.Logs)
		})
	}
}
