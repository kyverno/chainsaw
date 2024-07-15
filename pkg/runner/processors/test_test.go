package processors

// import (
// 	"context"
// 	"errors"
// 	"sync/atomic"
// 	"time"

// 	"github.com/jmespath-community/go-jmespath/pkg/binding"
// 	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
// 	"github.com/kyverno/chainsaw/pkg/client"
// 	fake "github.com/kyverno/chainsaw/pkg/client/testing"
// 	"github.com/kyverno/chainsaw/pkg/discovery"
// 	"github.com/kyverno/chainsaw/pkg/report"
// 	fakeNamespacer "github.com/kyverno/chainsaw/pkg/runner/namespacer/testing"
// 	"github.com/kyverno/chainsaw/pkg/runner/summary"
// 	"github.com/kyverno/chainsaw/pkg/testing"
// 	"github.com/stretchr/testify/assert"
// 	kerror "k8s.io/apimachinery/pkg/api/errors"
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/utils/clock"
// 	tclock "k8s.io/utils/clock/testing"
// 	"k8s.io/utils/ptr"
// 	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
// )

// func TestTestProcessor_Run(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		config         v1alpha1.ConfigurationSpec
// 		client         client.Client
// 		clock          clock.PassiveClock
// 		summary        *summary.Summary
// 		testsReport    *report.TestReport
// 		test           discovery.Test
// 		shouldFailFast bool
// 		binding        binding.Bindings
// 		namespacer     *fakeNamespacer.FakeNamespacer
// 		expectedFail   bool
// 		skipped        bool
// 	}{{
// 		name: "test with no steps",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Namespace: "chainsaw",
// 					Timeouts:  &v1alpha1.Timeouts{},
// 				},
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		namespacer:     nil,
// 		expectedFail:   false,
// 	}, {
// 		name: "test with test steps",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		namespacer:     nil,
// 		expectedFail:   false,
// 	}, {
// 		name: "fail fast",
// 		config: v1alpha1.ConfigurationSpec{
// 			FailFast: true,
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		shouldFailFast: true,
// 		binding:        binding.NewBindings(),
// 		namespacer:     nil,
// 		expectedFail:   false,
// 	}, {
// 		name: "skip test",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 					Skip:     ptr.To[bool](true),
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		expectedFail: false,
// 		skipped:      true,
// 	}, {
// 		name: "with test namespace",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Namespace: "chainsaw",
// 					Timeouts:  &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		expectedFail:   false,
// 	}, {
// 		name: "without test namespace",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Timeouts: &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		namespacer:     nil,
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		expectedFail:   false,
// 	}, {
// 		name: "delay before cleanup",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					DelayBeforeCleanup: ptr.To[v1.Duration](v1.Duration{Duration: 1 * time.Second}),
// 					Timeouts:           &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		namespacer:     nil,
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		expectedFail:   false,
// 	}, {
// 		name: "namespace not found and created",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("namespace"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Namespace: "chainsaw",
// 					Timeouts:  &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		expectedFail:   false,
// 	}, {
// 		name: "namespace not found due to internal error",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewInternalError(errors.New("internal error"))
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return nil
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Namespace: "chainsaw",
// 					Timeouts:  &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		expectedFail:   true,
// 	}, {
// 		name: "namespace not found and not created due to internal error",
// 		config: v1alpha1.ConfigurationSpec{
// 			Timeouts: v1alpha1.Timeouts{},
// 		},
// 		client: &fake.FakeClient{
// 			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
// 				return kerror.NewNotFound(v1alpha1.Resource("namespace"), "chainsaw")
// 			},
// 			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
// 				return kerror.NewInternalError(errors.New("internal error"))
// 			},
// 		},
// 		clock:       tclock.NewFakePassiveClock(time.Now()),
// 		summary:     &summary.Summary{},
// 		testsReport: nil,
// 		test: discovery.Test{
// 			Err: nil,
// 			Test: &v1alpha1.Test{
// 				Spec: v1alpha1.TestSpec{
// 					Namespace: "chainsaw",
// 					Timeouts:  &v1alpha1.Timeouts{},
// 					Steps: []v1alpha1.TestStep{
// 						{
// 							TestStepSpec: v1alpha1.TestStepSpec{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		namespacer: &fakeNamespacer.FakeNamespacer{
// 			GetNamespaceFn: func(call int) string {
// 				return "chainsaw"
// 			},
// 		},
// 		shouldFailFast: false,
// 		binding:        binding.NewBindings(),
// 		expectedFail:   true,
// 	}}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			shouldFailVar := &atomic.Bool{}
// 			shouldFailVar.Store(tc.shouldFailFast)
// 			registry := registryMock{}
// 			if tc.client != nil {
// 				registry.client = tc.client
// 			}
// 			processor := NewTestProcessor(
// 				tc.config,
// 				registry,
// 				tc.clock,
// 				tc.summary,
// 				tc.testsReport,
// 				tc.test,
// 				shouldFailVar,
// 			)
// 			nt := &testing.MockT{}
// 			ctx := testing.IntoContext(context.Background(), nt)
// 			processor.Run(ctx, tc.binding, tc.namespacer)
// 			nt.Cleanup(func() {})
// 			if tc.expectedFail {
// 				assert.True(t, nt.FailedVar, "expected an error but got none")
// 			} else {
// 				assert.False(t, nt.FailedVar, "expected no error but got one")
// 			}
// 			if shouldFailVar.Load() || tc.skipped {
// 				assert.True(t, nt.SkippedVar, "test should be skipped but it was not")
// 			} else {
// 				assert.False(t, nt.SkippedVar, "test should not be skipped but it was")
// 			}
// 		})
// 	}
// }
