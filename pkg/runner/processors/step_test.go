package processors

// import (
// 	"context"
// 	"path/filepath"
// 	"time"

// 	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
// 	"github.com/kyverno/chainsaw/pkg/client"
// 	fake "github.com/kyverno/chainsaw/pkg/client/testing"
// 	"github.com/kyverno/chainsaw/pkg/discovery"
// 	"github.com/kyverno/chainsaw/pkg/report"
// 	"github.com/kyverno/chainsaw/pkg/runner/logging"
// 	fakeLogger "github.com/kyverno/chainsaw/pkg/runner/logging/testing"
// 	fakeNamespacer "github.com/kyverno/chainsaw/pkg/runner/namespacer/testing"
// 	"github.com/kyverno/chainsaw/pkg/testing"
// 	"github.com/stretchr/testify/assert"
// 	kerror "k8s.io/apimachinery/pkg/api/errors"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/utils/clock"
// 	tclock "k8s.io/utils/clock/testing"
// 	"k8s.io/utils/ptr"
// 	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
// )

// func TestStepProcessor_Run(t *testing.T) {
// 	testData := filepath.Join("..", "..", "..", "testdata", "runner", "processors")
// 	testCases := []struct {
// 		name         string
// 		config       v1alpha1.ConfigurationSpec
// 		client       client.Client
// 		namespacer   *fakeNamespacer.FakeNamespacer
// 		clock        clock.PassiveClock
// 		test         discovery.Test
// 		stepSpec     v1alpha1.TestStep
// 		stepReport   *report.StepReport
// 		expectedFail bool
// 		skipped      bool
// 	}{{
// 		name: "test with no handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try:      []v1alpha1.Operation{},
// 				Catch:    []v1alpha1.CatchFinally{},
// 				Finally:  []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "try operation with apply handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 			PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Apply: &v1alpha1.Apply{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with create handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with assert handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				obj.(*unstructured.Unstructured).Object = map[string]any{
// 					"apiVersion": "v1",
// 					"kind":       "Pod",
// 					"metadata": map[string]any{
// 						"name": "myapp",
// 						"labels": map[string]string{
// 							"name": "myapp",
// 						},
// 					},
// 					"spec": map[string]any{
// 						"containers": []map[string]any{
// 							{
// 								"name":  "myapp",
// 								"image": "myapp:latest",
// 								"resources": map[string]any{
// 									"limits": map[string]string{
// 										"memory": "128Mi",
// 										"cpu":    "500m",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				}
// 				return nil
// 			},
// 			ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Assert: &v1alpha1.Assert{
// 							ActionCheckRef: v1alpha1.ActionCheckRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with error handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				obj.(*unstructured.Unstructured).Object = map[string]any{
// 					"apiVersion": "v1",
// 					"kind":       "Pod",
// 					"metadata": map[string]any{
// 						"name": "myapp",
// 						"labels": map[string]string{
// 							"name": "myapp",
// 						},
// 					},
// 					"spec": map[string]any{
// 						"containers": []map[string]any{
// 							{
// 								"name":  "myapp",
// 								"image": "myapp:fake",
// 								"resources": map[string]any{
// 									"limits": map[string]string{
// 										"memory": "128Mi",
// 										"cpu":    "500m",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				}
// 				return nil
// 			},
// 			ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Error: &v1alpha1.Error{
// 							ActionCheckRef: v1alpha1.ActionCheckRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with command handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Command: &v1alpha1.Command{
// 							Entrypoint: "echo",
// 							Args:       []string{"hello"},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with script handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Script: &v1alpha1.Script{
// 							Content: "echo hello",
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with sleep handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client:     &fake.FakeClient{},
// 		namespacer: &fakeNamespacer.FakeNamespacer{},
// 		clock:      tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Sleep: &v1alpha1.Sleep{
// 							Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "try operation with delete handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("Deployment"), "chainsaw")
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Delete: &v1alpha1.Delete{
// 							Ref: &v1alpha1.ObjectReference{
// 								ObjectType: v1alpha1.ObjectType{
// 									APIVersion: "apps/v1",
// 									Kind:       "Deployment",
// 								},
// 								ObjectName: v1alpha1.ObjectName{
// 									Namespace: "chainsaw",
// 									Name:      "myapp",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "dry run with create handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 							ActionDryRun: v1alpha1.ActionDryRun{
// 								DryRun: ptr.To[bool](true),
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "skip delete with create handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				SkipDelete: ptr.To[bool](true),
// 				Timeouts:   &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "try-raw resource with create handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				SkipDelete: ptr.To[bool](true),
// 				Timeouts:   &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								Resource: &unstructured.Unstructured{
// 									Object: map[string]any{
// 										"apiVersion": "v1",
// 										"kind":       "Pod",
// 										"metadata": map[string]any{
// 											"name": "chainsaw",
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "try-url resource with create handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("pod"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				SkipDelete: ptr.To[bool](true),
// 				Timeouts:   &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/test/configmap.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "raw resource with assert handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				obj.(*unstructured.Unstructured).Object = map[string]any{
// 					"apiVersion": "v1",
// 					"kind":       "Pod",
// 					"metadata": map[string]any{
// 						"name": "myapp",
// 						"labels": map[string]string{
// 							"name": "myapp",
// 						},
// 					},
// 					"spec": map[string]any{
// 						"containers": []map[string]any{
// 							{
// 								"name":  "myapp",
// 								"image": "myapp:latest",
// 								"resources": map[string]any{
// 									"limits": map[string]string{
// 										"memory": "128Mi",
// 										"cpu":    "500m",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				}
// 				return nil
// 			},
// 			ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Assert: &v1alpha1.Assert{
// 							ActionCheckRef: v1alpha1.ActionCheckRef{
// 								Check: &v1alpha1.Check{
// 									Value: map[string]any{
// 										"apiVersion": "v1",
// 										"kind":       "Pod",
// 										"metadata": map[string]any{
// 											"name": "myapp",
// 											"labels": map[string]string{
// 												"name": "myapp",
// 											},
// 										},
// 										"spec": map[string]any{
// 											"containers": []map[string]any{
// 												{
// 													"name":  "myapp",
// 													"image": "myapp:latest",
// 													"resources": map[string]any{
// 														"limits": map[string]string{
// 															"memory": "128Mi",
// 															"cpu":    "500m",
// 														},
// 													},
// 												},
// 											},
// 										},
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "try url-resource with assert handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				obj.(*unstructured.Unstructured).Object = map[string]any{
// 					"apiVersion": "v1",
// 					"kind":       "ConfigMap",
// 					"metadata": map[string]any{
// 						"name": "default",
// 					},
// 					"data": map[string]string{
// 						"foo": "bar",
// 					},
// 				}
// 				return nil
// 			},
// 			ListFn: func(ctx context.Context, call int, list ctrlclient.ObjectList, opts ...ctrlclient.ListOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Assert: &v1alpha1.Assert{
// 							ActionCheckRef: v1alpha1.ActionCheckRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "https://raw.githubusercontent.com/kyverno/chainsaw/main/testdata/test/configmap.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}, {
// 		name: "try, catch and finally operation with apply handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			ForceTerminationGracePeriod: &metav1.Duration{Duration: time.Duration(1) * time.Second},
// 			Timeouts:                    v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				obj.(*unstructured.Unstructured).Object = map[string]any{
// 					"apiVersion": "v1",
// 					"kind":       "Pod",
// 					"metadata": map[string]any{
// 						"name":      "myapp",
// 						"namespace": "chainsaw",
// 						"labels": map[string]string{
// 							"name": "myapp",
// 						},
// 					},
// 					"spec": map[string]any{
// 						"containers": []map[string]any{
// 							{
// 								"name":  "myapp",
// 								"image": "myapp:latest",
// 								"resources": map[string]any{
// 									"limits": map[string]string{
// 										"memory": "128Mi",
// 										"cpu":    "500m",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				}
// 				return nil
// 			},
// 			PatchFn: func(ctx context.Context, call int, obj ctrlclient.Object, patch ctrlclient.Patch, opts ...ctrlclient.PatchOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Apply: &v1alpha1.Apply{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "pod.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch: []v1alpha1.CatchFinally{
// 					{
// 						Command: &v1alpha1.Command{
// 							Entrypoint: "echo",
// 							Args:       []string{"hello"},
// 						},
// 					},
// 					{
// 						Script: &v1alpha1.Script{
// 							Content: "echo hello",
// 						},
// 					},
// 					{
// 						Sleep: &v1alpha1.Sleep{
// 							Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
// 						},
// 					},
// 					{
// 						PodLogs: &v1alpha1.PodLogs{
// 							ActionObjectSelector: v1alpha1.ActionObjectSelector{
// 								Selector: "name=myapp",
// 							},
// 						},
// 					},
// 				},
// 				Finally: []v1alpha1.CatchFinally{
// 					{
// 						Command: &v1alpha1.Command{
// 							Entrypoint: "echo",
// 							Args:       []string{"hello"},
// 						},
// 					},
// 					{
// 						Script: &v1alpha1.Script{
// 							Content: "echo hello",
// 						},
// 					},
// 					{
// 						Sleep: &v1alpha1.Sleep{
// 							Duration: metav1.Duration{Duration: time.Duration(1) * time.Second},
// 						},
// 					},
// 					{
// 						PodLogs: &v1alpha1.PodLogs{
// 							ActionObjectSelector: v1alpha1.ActionObjectSelector{
// 								Selector: "name=myapp",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		stepReport: &report.StepReport{},
// 	}, {
// 		name: "termination with create handler",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("Deployment"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			ApplyFn: func(obj ctrlclient.Object, call int) error {
// 				return nil
// 			},
// 		},
// 		clock: tclock.NewFakePassiveClock(time.Now()),
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					ForceTerminationGracePeriod: &metav1.Duration{Duration: time.Duration(1) * time.Second},
// 					Timeouts:                    &v1alpha1.Timeouts{},
// 				},
// 			},
// 			BasePath: testData,
// 		},
// 		stepSpec: v1alpha1.TestStep{
// 			TestStepSpec: v1alpha1.TestStepSpec{
// 				Timeouts: &v1alpha1.Timeouts{},
// 				Try: []v1alpha1.Operation{
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "deployment.yaml",
// 								},
// 							},
// 						},
// 					},
// 					{
// 						Create: &v1alpha1.Create{
// 							ActionResourceRef: v1alpha1.ActionResourceRef{
// 								FileRef: v1alpha1.FileRef{
// 									File: "cron-job.yaml",
// 								},
// 							},
// 						},
// 					},
// 				},
// 				Catch:   []v1alpha1.CatchFinally{},
// 				Finally: []v1alpha1.CatchFinally{},
// 			},
// 		},
// 		stepReport: nil,
// 	}}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			registry := registryMock{}
// 			if tc.client != nil {
// 				registry.client = tc.client
// 			}
// 			stepProcessor := NewStepProcessor(
// 				tc.config,
// 				registry,
// 				tc.namespacer,
// 				tc.clock,
// 				tc.test,
// 				tc.stepSpec,
// 				tc.stepReport,
// 			)
// 			nt := &testing.MockT{}
// 			ctx := testing.IntoContext(context.Background(), nt)
// 			ctx = logging.IntoContext(ctx, &fakeLogger.FakeLogger{})
// 			stepProcessor.Run(ctx, nil)
// 			nt.Cleanup(func() {})
// 			if tc.expectedFail {
// 				assert.True(t, nt.FailedVar, "expected an error but got none")
// 			} else {
// 				assert.False(t, nt.FailedVar, "expected no error but got one")
// 			}
// 		})
// 	}
// }
