package runner

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	fake "github.com/kyverno/chainsaw/pkg/client/testing"
	fakeNamespacer "github.com/kyverno/chainsaw/pkg/engine/namespacer/testing"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/logging"
	fakeLogger "github.com/kyverno/chainsaw/pkg/mocks"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/mocks"
	"github.com/stretchr/testify/assert"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/clock"
	"k8s.io/utils/ptr"
)

func TestStepProcessor_Run(t *testing.T) {
	config, err := config.DefaultConfiguration()
	if err != nil {
		assert.NoError(t, err)
	}
	testData := filepath.Join("..", "..", "testdata", "runner", "processors")
	testCases := []struct {
		name                   string
		client                 client.Client
		namespacer             *fakeNamespacer.FakeNamespacer
		basePath               string
		terminationGracePeriod *metav1.Duration
		stepSpec               v1alpha1.TestStep
		want                   bool
		expectedFail           bool
	}{{
		name:   "test with no handler",
		client: &fake.FakeClient{},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: "",
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try:      []v1alpha1.Operation{},
				Catch:    []v1alpha1.CatchFinally{},
				Finally:  []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try operation with apply handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return nil
			},
			PatchFn: func(ctx context.Context, call int, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Apply: &v1alpha1.Apply{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try operation with create handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try operation with assert handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "myapp",
						"labels": map[string]string{
							"name": "myapp",
						},
					},
					"spec": map[string]any{
						"containers": []map[string]any{
							{
								"name":  "myapp",
								"image": "myapp:latest",
								"resources": map[string]any{
									"limits": map[string]string{
										"memory": "128Mi",
										"cpu":    "500m",
									},
								},
							},
						},
					},
				}
				return nil
			},
			ListFn: func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Assert: &v1alpha1.Assert{
							ActionCheckRef: v1alpha1.ActionCheckRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try operation with error handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "myapp",
						"labels": map[string]string{
							"name": "myapp",
						},
					},
					"spec": map[string]any{
						"containers": []map[string]any{
							{
								"name":  "myapp",
								"image": "myapp:fake",
								"resources": map[string]any{
									"limits": map[string]string{
										"memory": "128Mi",
										"cpu":    "500m",
									},
								},
							},
						},
					},
				}
				return nil
			},
			ListFn: func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Error: &v1alpha1.Error{
							ActionCheckRef: v1alpha1.ActionCheckRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name:   "try operation with command handler",
		client: &fake.FakeClient{},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Command: &v1alpha1.Command{
							Entrypoint: "echo",
							Args:       []string{"hello"},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name:   "try operation with script handler",
		client: &fake.FakeClient{},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Script: &v1alpha1.Script{
							Content: "echo hello",
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name:       "try operation with sleep handler",
		client:     &fake.FakeClient{},
		namespacer: &fakeNamespacer.FakeNamespacer{},
		basePath:   testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Sleep: &v1alpha1.Sleep{
							Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try operation with delete handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("Deployment"), "chainsaw")
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Delete: &v1alpha1.Delete{
							Ref: &v1alpha1.ObjectReference{
								ObjectType: v1alpha1.ObjectType{
									APIVersion: "apps/v1",
									Kind:       "Deployment",
								},
								ObjectName: v1alpha1.ObjectName{
									Namespace: "chainsaw",
									Name:      "myapp",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "dry run with create handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
							ActionDryRun: v1alpha1.ActionDryRun{
								DryRun: ptr.To[bool](true),
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "skip delete with create handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				SkipDelete: ptr.To[bool](true),
				Timeouts:   &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try-raw resource with create handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				SkipDelete: ptr.To[bool](true),
				Timeouts:   &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								Resource: &unstructured.Unstructured{
									Object: map[string]any{
										"apiVersion": "v1",
										"kind":       "Pod",
										"metadata": map[string]any{
											"name": "chainsaw",
										},
									},
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try-url resource with create handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				SkipDelete: ptr.To[bool](true),
				Timeouts:   &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/test/configmap.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "raw resource with assert handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name": "myapp",
						"labels": map[string]string{
							"name": "myapp",
						},
					},
					"spec": map[string]any{
						"containers": []map[string]any{
							{
								"name":  "myapp",
								"image": "myapp:latest",
								"resources": map[string]any{
									"limits": map[string]string{
										"memory": "128Mi",
										"cpu":    "500m",
									},
								},
							},
						},
					},
				}
				return nil
			},
			ListFn: func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Assert: &v1alpha1.Assert{
							ActionCheckRef: v1alpha1.ActionCheckRef{
								Check: ptr.To(v1alpha1.NewProjection(
									map[string]any{
										"apiVersion": "v1",
										"kind":       "Pod",
										"metadata": map[string]any{
											"name": "myapp",
											"labels": map[string]string{
												"name": "myapp",
											},
										},
										"spec": map[string]any{
											"containers": []map[string]any{
												{
													"name":  "myapp",
													"image": "myapp:latest",
													"resources": map[string]any{
														"limits": map[string]string{
															"memory": "128Mi",
															"cpu":    "500m",
														},
													},
												},
											},
										},
									},
								)),
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try url-resource with assert handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]any{
						"name": "default",
					},
					"data": map[string]string{
						"foo": "bar",
					},
				}
				return nil
			},
			ListFn: func(ctx context.Context, call int, list client.ObjectList, opts ...client.ListOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath: testData,
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Assert: &v1alpha1.Assert{
							ActionCheckRef: v1alpha1.ActionCheckRef{
								FileRef: v1alpha1.FileRef{
									File: "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/test/configmap.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}, {
		name: "try, catch and finally operation with apply handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				obj.(*unstructured.Unstructured).Object = map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name":      "myapp",
						"namespace": "chainsaw",
						"labels": map[string]string{
							"name": "myapp",
						},
					},
					"spec": map[string]any{
						"containers": []map[string]any{
							{
								"name":  "myapp",
								"image": "myapp:latest",
								"resources": map[string]any{
									"limits": map[string]string{
										"memory": "128Mi",
										"cpu":    "500m",
									},
								},
							},
						},
					},
				}
				return nil
			},
			PatchFn: func(ctx context.Context, call int, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		basePath:               testData,
		terminationGracePeriod: &metav1.Duration{Duration: time.Duration(1) * time.Second},
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Apply: &v1alpha1.Apply{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "pod.yaml",
								},
							},
						},
					},
				},
				Catch: []v1alpha1.CatchFinally{
					{
						Command: &v1alpha1.Command{
							Entrypoint: "echo",
							Args:       []string{"hello"},
						},
					},
					{
						Script: &v1alpha1.Script{
							Content: "echo hello",
						},
					},
					{
						Sleep: &v1alpha1.Sleep{
							Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
						},
					},
					{
						PodLogs: &v1alpha1.PodLogs{
							ActionObjectSelector: v1alpha1.ActionObjectSelector{
								Selector: "name=myapp",
							},
						},
					},
				},
				Finally: []v1alpha1.CatchFinally{
					{
						Command: &v1alpha1.Command{
							Entrypoint: "echo",
							Args:       []string{"hello"},
						},
					},
					{
						Script: &v1alpha1.Script{
							Content: "echo hello",
						},
					},
					{
						Sleep: &v1alpha1.Sleep{
							Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
						},
					},
					{
						PodLogs: &v1alpha1.PodLogs{
							ActionObjectSelector: v1alpha1.ActionObjectSelector{
								Selector: "name=myapp",
							},
						},
					},
				},
			},
		},
	}, {
		name: "termination with create handler",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("Deployment"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj client.Object, opts ...client.CreateOption) error {
				return nil
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			ApplyFn: func(int, client.Client, client.Object) error {
				return nil
			},
		},
		basePath:               testData,
		terminationGracePeriod: &metav1.Duration{Duration: time.Duration(1) * time.Second},
		stepSpec: v1alpha1.TestStep{
			TestStepSpec: v1alpha1.TestStepSpec{
				Timeouts: &v1alpha1.Timeouts{},
				Try: []v1alpha1.Operation{
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "deployment.yaml",
								},
							},
						},
					},
					{
						Create: &v1alpha1.Create{
							ActionResourceRef: v1alpha1.ActionResourceRef{
								FileRef: v1alpha1.FileRef{
									File: "cron-job.yaml",
								},
							},
						},
					},
				},
				Catch:   []v1alpha1.CatchFinally{},
				Finally: []v1alpha1.CatchFinally{},
			},
		},
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registry := mocks.Registry{}
			if tc.client != nil {
				registry.Client = tc.client
			}
			_failed := false
			fail := func() { _failed = true }
			failed := func() bool { return _failed }
			cleanup := func(func()) {}
			ctx := context.Background()
			ctx = logging.WithLogger(ctx, &fakeLogger.Logger{})
			tcontext := enginecontext.MakeContext(clock.RealClock{}, apis.NewBindings(), registry).WithTimeouts(v1alpha1.Timeouts{
				Apply:   &config.Spec.Timeouts.Apply,
				Assert:  &config.Spec.Timeouts.Assert,
				Cleanup: &config.Spec.Timeouts.Cleanup,
				Delete:  &config.Spec.Timeouts.Delete,
				Error:   &config.Spec.Timeouts.Error,
				Exec:    &config.Spec.Timeouts.Exec,
			})
			runner := runner{}
			got := runner.runStep(ctx, cleanup, fail, failed, tc.basePath, tc.namespacer, tcontext, tc.stepSpec, &model.TestReport{})
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.expectedFail, _failed)
		})
	}
}
