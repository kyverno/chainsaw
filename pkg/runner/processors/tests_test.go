package processors

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	fake "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/summary"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/clock"
	tclock "k8s.io/utils/clock/testing"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestTestsProcessor_Run(t *testing.T) {
	testCases := []struct {
		name         string
		config       v1alpha1.ConfigurationSpec
		client       client.Client
		clock        clock.PassiveClock
		summary      *summary.Summary
		testsReport  *report.TestsReport
		bindings     binding.Bindings
		tests        []discovery.Test
		expectedFail bool
	}{
		{
			name: "Namesapce exists",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "default",
			},
			client: &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
			},
			clock:        nil,
			summary:      &summary.Summary{},
			testsReport:  &report.TestsReport{},
			bindings:     binding.NewBindings(),
			tests:        []discovery.Test{},
			expectedFail: false,
		},
		{
			name: "Namesapce doesn't exists",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "chain-saw",
			},
			client: &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return errors.NewNotFound(v1alpha1.Resource("Namespace"), "chain-saw")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return nil
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					return nil
				},
			},
			clock:        nil,
			summary:      &summary.Summary{},
			testsReport:  &report.TestsReport{},
			bindings:     binding.NewBindings(),
			tests:        []discovery.Test{},
			expectedFail: false,
		},
		{
			name: "Namesapce not found with error",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "chain-saw",
			},
			client: &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return errors.NewBadRequest("failed to get namespace")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return nil
				},
			},
			clock:        nil,
			summary:      &summary.Summary{},
			testsReport:  &report.TestsReport{},
			bindings:     binding.NewBindings(),
			tests:        []discovery.Test{},
			expectedFail: true,
		},
		{
			name: "Namesapce doesn't exists and can't be created",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "chain-saw",
			},
			client: &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return errors.NewNotFound(v1alpha1.Resource("Namespace"), "chain-saw")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return errors.NewBadRequest("failed to create namespace")
				},
			},
			clock:        nil,
			summary:      &summary.Summary{},
			testsReport:  &report.TestsReport{},
			bindings:     binding.NewBindings(),
			tests:        []discovery.Test{},
			expectedFail: true,
		},
		{
			name: "Success",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "default",
			},
			client: &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
			},
			clock:       nil,
			summary:     &summary.Summary{},
			testsReport: &report.TestsReport{},
			bindings:    binding.NewBindings(),
			tests: []discovery.Test{
				{
					Err:      nil,
					BasePath: "fakePath",
					Test:     &v1alpha1.Test{},
				},
			},
			expectedFail: false,
		},
		{
			name: "Fail",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "default",
			},
			client: &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
			},
			clock:       nil,
			summary:     &summary.Summary{},
			testsReport: &report.TestsReport{},
			bindings:    binding.NewBindings(),
			tests: []discovery.Test{
				{
					Err:      errors.NewBadRequest("failed to get test"),
					BasePath: "fakePath",
					Test:     nil,
				},
			},
			expectedFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clusters := NewClusters()
			if tc.client != nil {
				clusters.clients[defaultClient] = tc.client
			}
			processor := NewTestsProcessor(
				tc.config,
				clusters,
				tc.clock,
				tc.summary,
				tc.testsReport,
				tc.bindings,
				tc.tests...,
			)
			nt := testing.MockT{}
			ctx := testing.IntoContext(context.Background(), &nt)
			processor.Run(ctx)
			nt.Cleanup(func() {
			})
			if tc.expectedFail {
				assert.True(t, nt.FailedVar, "expected an error but got none")
			} else {
				assert.False(t, nt.FailedVar, "expected no error but got one")
			}
		})
	}
}

func TestCreateTestProcessor(t *testing.T) {
	testCases := []struct {
		name        string
		config      v1alpha1.ConfigurationSpec
		client      client.Client
		clock       clock.PassiveClock
		summary     *summary.Summary
		testsReport *report.TestsReport
		test        []discovery.Test
	}{
		{
			name: "TestProcessor is created",
			config: v1alpha1.ConfigurationSpec{
				Namespace: "default",
			},
			client:      &fake.FakeClient{},
			clock:       tclock.NewFakePassiveClock(time.Now()),
			summary:     &summary.Summary{},
			testsReport: report.NewTests("FakeReport"),
			test: []discovery.Test{
				{
					Err:      nil,
					BasePath: "fakePath",
					Test:     &v1alpha1.Test{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			localTC := tc
			clusters := NewClusters()
			if localTC.client != nil {
				clusters.clients[defaultClient] = localTC.client
			}
			processor := testsProcessor{
				config:         localTC.config,
				clusters:       clusters,
				clock:          localTC.clock,
				summary:        localTC.summary,
				testsReport:    localTC.testsReport,
				tests:          localTC.test,
				shouldFailFast: atomic.Bool{},
			}
			processor.shouldFailFast.Store(false)

			result := processor.CreateTestProcessor(localTC.test[0], nil)

			assert.NotNil(t, result, "TestProcessor should not be nil")
			if localTC.testsReport != nil {
				assert.True(t, len(localTC.testsReport.Reports) > 0, "Test report should be added")
			}
		})
	}
}
