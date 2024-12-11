package runner

import (
	"context"
	"errors"
	"flag"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/flags"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
	tclock "k8s.io/utils/clock/testing"
)

func TestNew(t *testing.T) {
	realClock := clock.RealClock{}
	onFailure := func() {}
	tests := []struct {
		name      string
		clock     clock.PassiveClock
		onFailure func()
	}{{
		clock: realClock,
	}, {
		clock:     realClock,
		onFailure: onFailure,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.clock, tt.onFailure)
			runner, ok := got.(*runner)
			assert.True(t, ok)
			assert.Equal(t, tt.clock, runner.clock)
			if tt.onFailure == nil {
				assert.Nil(t, runner.onFailure)
			} else {
				assert.NotNil(t, runner.onFailure)
			}
		})
	}
}

func Test_runner_Run(t *testing.T) {
	config, err := config.DefaultConfiguration()
	assert.NoError(t, err)
	tc, err := InitContext(config.Spec, nil, nil)
	assert.NoError(t, err)
	type summaryResult struct {
		passed  int32
		failed  int32
		skipped int32
	}
	tests := []struct {
		name    string
		config  model.Configuration
		tc      enginecontext.TestContext
		tests   []discovery.Test
		want    *summaryResult
		wantErr bool
	}{{
		name:    "no tests",
		config:  config.Spec,
		tc:      tc,
		tests:   nil,
		want:    nil,
		wantErr: false,
	}, {
		name:   "test with err",
		config: config.Spec,
		tc:     tc,
		tests:  []discovery.Test{{Err: errors.New("dummy")}},
		want: &summaryResult{
			passed:  0,
			failed:  1,
			skipped: 0,
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &runner{
				clock: clock.RealClock{},
				deps:  &internal.TestDeps{Test: true},
			}
			// to allow unit tests to work
			assert.NoError(t, flags.SetupFlags(tt.config))
			assert.NoError(t, flag.Set("test.testlogfile", ""))
			got, err := r.Run(context.TODO(), tt.config, tt.tc, tt.tests...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.failed, got.Failed())
				assert.Equal(t, tt.want.passed, got.Passed())
				assert.Equal(t, tt.want.skipped, got.Skipped())
			}
		})
	}
}

type _mainStart struct {
	code int
}

func (m *_mainStart) Run() int {
	return m.code
}

func Test_runner_run(t *testing.T) {
	runTests := []discovery.Test{{}}
	tests := []struct {
		name    string
		m       mainstart
		tests   []discovery.Test
		wantErr bool
	}{{
		name:    "with 0",
		m:       &_mainStart{code: 0},
		tests:   runTests,
		wantErr: false,
	}, {
		name:    "with 1",
		m:       &_mainStart{code: 1},
		tests:   runTests,
		wantErr: false,
	}, {
		name:    "with 2",
		m:       &_mainStart{code: 2},
		tests:   runTests,
		wantErr: true,
	}, {
		name:    "no tests",
		m:       &_mainStart{code: 2},
		tests:   nil,
		wantErr: false,
	}, {
		name:    "no tests",
		m:       &_mainStart{code: 2},
		tests:   []discovery.Test{},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &runner{
				clock: clock.RealClock{},
			}
			_, err := r.run(context.TODO(), tt.m, model.Configuration{}, enginecontext.EmptyContext(), tt.tests...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_runner_onFail(t *testing.T) {
	{
		r := &runner{
			clock: clock.RealClock{},
		}
		r.onFail()
	}
	{
		called := false
		r := &runner{
			clock: clock.RealClock{},
			onFailure: func() {
				called = true
			},
		}
		r.onFail()
		assert.True(t, called)
	}
}

func TestRun(t *testing.T) {
	fakeClock := tclock.NewFakePassiveClock(time.Now())
	config, err := config.DefaultConfiguration()
	if err != nil {
		assert.NoError(t, err)
	}
	tests := []struct {
		name       string
		tests      []discovery.Test
		config     model.Configuration
		restConfig *rest.Config
		mockReturn int
		wantErr    bool
	}{{
		name:  "Zero Tests",
		tests: []discovery.Test{},
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
		},
		restConfig: &rest.Config{},
		wantErr:    false,
	}, {
		name: "Nil Rest Config with 1 Test",
		tests: []discovery.Test{
			{
				Err: nil,
				Test: &model.Test{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
			Report: &v1alpha2.ReportOptions{
				Format: v1alpha2.JSONFormat,
			},
		},
		restConfig: nil,
		wantErr:    false,
	}, {
		name:  "Zero Tests with JSON Report",
		tests: []discovery.Test{},
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
			Report: &v1alpha2.ReportOptions{
				Format: v1alpha2.JSONFormat,
			},
		},
		restConfig: &rest.Config{},
		wantErr:    false,
	}, {
		name: "Success Case with 1 Test",
		tests: []discovery.Test{
			{
				Err: nil,
				Test: &model.Test{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
		restConfig: &rest.Config{},
		mockReturn: 0,
		wantErr:    false,
	}, {
		name: "Failure Case with 1 Test",
		tests: []discovery.Test{
			{
				Err: nil,
				Test: &model.Test{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
		restConfig: &rest.Config{},
		mockReturn: 2,
		wantErr:    true,
	}, {
		name: "Success Case with 1 Test with XML Report",
		tests: []discovery.Test{
			{
				Err: nil,
				Test: &model.Test{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
			Report: &v1alpha2.ReportOptions{
				Format: v1alpha2.XMLFormat,
				Name:   "chainsaw",
			},
		},
		restConfig: &rest.Config{},
		mockReturn: 0,
		wantErr:    false,
	}, {
		name: "Error in saving Report",
		tests: []discovery.Test{
			{
				Err: nil,
				Test: &model.Test{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
				},
			},
		},
		config: model.Configuration{
			Timeouts: config.Spec.Timeouts,
			Report: &v1alpha2.ReportOptions{
				Format: "abc",
			},
		},
		restConfig: &rest.Config{},
		mockReturn: 0,
		wantErr:    true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMainStart := &_mainStart{
				code: tt.mockReturn,
			}
			runner := runner{
				clock: fakeClock,
			}
			ctx := context.TODO()
			tc, err := InitContext(tt.config, tt.restConfig, nil)
			assert.NoError(t, err)
			_, err = runner.run(ctx, mockMainStart, tt.config, tc, tt.tests...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
