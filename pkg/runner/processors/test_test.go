package processors

import (
	"context"
	"errors"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/client"
	fake "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/discovery"
	enginecontext "github.com/kyverno/chainsaw/pkg/engine/context"
	fakeNamespacer "github.com/kyverno/chainsaw/pkg/engine/namespacer/testing"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	tclock "k8s.io/utils/clock/testing"
	"k8s.io/utils/ptr"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestTestProcessor_Run(t *testing.T) {
	config, err := config.DefaultConfiguration()
	if err != nil {
		assert.NoError(t, err)
	}
	testCases := []struct {
		name           string
		config         model.Configuration
		client         client.Client
		clock          clock.PassiveClock
		testsReport    *report.TestReport
		test           discovery.Test
		shouldFailFast bool
		binding        binding.Bindings
		namespacer     *fakeNamespacer.FakeNamespacer
		expectedFail   bool
		skipped        bool
	}{{
		name: "test with no steps",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Namespace: "chainsaw",
					Timeouts:  &v1alpha1.Timeouts{},
				},
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		namespacer:     nil,
		expectedFail:   false,
	}, {
		name: "test with test steps",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Timeouts: &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		namespacer:     nil,
		expectedFail:   false,
	}, {
		name: "fail fast",
		config: model.Configuration{
			Execution: v1alpha2.ExecutionOptions{
				FailFast: true,
			},
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Timeouts: &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		shouldFailFast: true,
		binding:        binding.NewBindings(),
		namespacer:     nil,
		expectedFail:   false,
	}, {
		name: "skip test",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Timeouts: &v1alpha1.Timeouts{},
					Skip:     ptr.To[bool](true),
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		namespacer: &fakeNamespacer.FakeNamespacer{
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		expectedFail: false,
		skipped:      true,
	}, {
		name: "with test namespace",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Namespace: "chainsaw",
					Timeouts:  &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		expectedFail:   false,
	}, {
		name: "without test namespace",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Timeouts: &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		namespacer:     nil,
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		expectedFail:   false,
	}, {
		name: "delay before cleanup",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					DelayBeforeCleanup: ptr.To[v1.Duration](v1.Duration{Duration: 1 * time.Second}),
					Timeouts:           &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		namespacer:     nil,
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		expectedFail:   false,
	}, {
		name: "namespace not found and created",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("namespace"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Namespace: "chainsaw",
					Timeouts:  &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		expectedFail:   false,
	}, {
		name: "namespace not found due to internal error",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerror.NewInternalError(errors.New("internal error"))
			},
			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Namespace: "chainsaw",
					Timeouts:  &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		expectedFail:   true,
	}, {
		name: "namespace not found and not created due to internal error",
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("namespace"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				return kerror.NewInternalError(errors.New("internal error"))
			},
		},
		clock:       tclock.NewFakePassiveClock(time.Now()),
		testsReport: nil,
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Namespace: "chainsaw",
					Timeouts:  &v1alpha1.Timeouts{},
					Steps: []v1alpha1.TestStep{
						{
							TestStepSpec: v1alpha1.TestStepSpec{},
						},
					},
				},
			},
		},
		namespacer: &fakeNamespacer.FakeNamespacer{
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		shouldFailFast: false,
		binding:        binding.NewBindings(),
		expectedFail:   true,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registry := registryMock{}
			if tc.client != nil {
				registry.client = tc.client
			}
			processor := NewTestProcessor(tc.config, tc.clock, tc.testsReport, 0)
			nt := &testing.MockT{}
			ctx := testing.IntoContext(context.Background(), nt)
			tcontext := enginecontext.MakeContext(binding.NewBindings(), registry)
			processor.Run(ctx, tc.namespacer, tcontext, tc.test)
			nt.Cleanup(func() {})
			if tc.expectedFail {
				assert.True(t, nt.FailedVar, "expected an error but got none")
			} else {
				assert.False(t, nt.FailedVar, "expected no error but got one")
			}
			// if shouldFailVar.Load() || tc.skipped {
			// 	assert.True(t, nt.SkippedVar, "test should be skipped but it was not")
			// } else {
			// 	assert.False(t, nt.SkippedVar, "test should not be skipped but it was")
			// }
		})
	}
}
