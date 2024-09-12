package processors

import (
	"context"
	"errors"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	fake "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/discovery"
	enginecontext "github.com/kyverno/chainsaw/pkg/engine/context"
	fakeNamespacer "github.com/kyverno/chainsaw/pkg/engine/namespacer/testing"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
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
		name         string
		client       client.Client
		clock        clock.PassiveClock
		test         discovery.Test
		namespacer   *fakeNamespacer.FakeNamespacer
		expectedFail bool
		skipped      bool
	}{{
		name: "test with no steps",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
		test: discovery.Test{
			Err: nil,
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Namespace: "chainsaw",
					Timeouts:  &v1alpha1.Timeouts{},
				},
			},
		},
		namespacer:   nil,
		expectedFail: false,
	}, {
		name: "test with test steps",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		namespacer:   nil,
		expectedFail: false,
	}, {
		name: "skip test",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		namespacer: &fakeNamespacer.FakeNamespacer{
			GetNamespaceFn: func(call int) string {
				return "chainsaw"
			},
		},
		expectedFail: false,
		skipped:      true,
	}, {
		name: "with test namespace",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		expectedFail: false,
	}, {
		name: "without test namespace",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		namespacer:   nil,
		expectedFail: false,
	}, {
		name: "delay before cleanup",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		namespacer:   nil,
		expectedFail: false,
	}, {
		name: "namespace not found and created",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("namespace"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		expectedFail: false,
	}, {
		name: "namespace not found due to internal error",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerror.NewInternalError(errors.New("internal error"))
			},
			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				return nil
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		expectedFail: true,
	}, {
		name: "namespace not found and not created due to internal error",
		client: &fake.FakeClient{
			GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
				return kerror.NewNotFound(v1alpha1.Resource("namespace"), "chainsaw")
			},
			CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
				return kerror.NewInternalError(errors.New("internal error"))
			},
		},
		clock: tclock.NewFakePassiveClock(time.Now()),
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
		expectedFail: true,
	}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registry := registryMock{}
			if tc.client != nil {
				registry.client = tc.client
			}
			processor := NewTestProcessor(
				tc.test,
				0,
				tc.clock,
				config.Spec.Namespace.Template,
				nil,
				config.Spec.Execution.ForceTerminationGracePeriod,
				config.Spec.Timeouts,
				config.Spec.Deletion.Propagation,
				config.Spec.Templating.Enabled,
				config.Spec.Cleanup.SkipDelete,
				config.Spec.Error.Catch...,
			)
			nt := &testing.MockT{}
			ctx := testing.IntoContext(context.Background(), nt)
			tcontext := enginecontext.MakeContext(binding.NewBindings(), registry)
			processor.Run(ctx, tc.namespacer, tcontext)
			if tc.expectedFail {
				assert.True(t, nt.FailedVar, "expected an error but got none")
			} else {
				assert.False(t, nt.FailedVar, "expected no error but got one")
			}
		})
	}
}
