package runner

import (
	"context"
	"errors"
	"flag"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha2"
	"github.com/kyverno/chainsaw/pkg/client"
	fake "github.com/kyverno/chainsaw/pkg/client/testing"
	"github.com/kyverno/chainsaw/pkg/discovery"
	"github.com/kyverno/chainsaw/pkg/loaders/config"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/chainsaw/pkg/runner/flags"
	"github.com/kyverno/chainsaw/pkg/runner/internal"
	"github.com/kyverno/chainsaw/pkg/runner/mocks"
	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/clock"
	"k8s.io/utils/ptr"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
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
	tc := func() enginecontext.TestContext {
		tc, err := enginecontext.InitContext(config.Spec, nil, nil)
		assert.NoError(t, err)
		return tc
	}
	echoTest := discovery.Test{
		Test: &model.Test{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.TestSpec{
				Steps: []v1alpha1.TestStep{{
					TestStepSpec: v1alpha1.TestStepSpec{
						Try: []v1alpha1.Operation{{
							Script: &v1alpha1.Script{
								Content: "echo hello",
							},
						}},
					},
				}},
			},
		},
	}
	scenarioTest := discovery.Test{
		Test: &model.Test{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "chainsaw.kyverno.io/v1alpha1",
				Kind:       "Test",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1alpha1.TestSpec{
				Scenarios: []v1alpha1.Scenario{{
					Bindings: []v1alpha1.Binding{{
						Name:  "foo",
						Value: v1alpha1.NewProjection("bar"),
					}},
				}},
				Steps: []v1alpha1.TestStep{{
					TestStepSpec: v1alpha1.TestStepSpec{
						Try: []v1alpha1.Operation{{
							Script: &v1alpha1.Script{
								ActionEnv: v1alpha1.ActionEnv{
									Env: []v1alpha1.Binding{{
										Name:  "FOO",
										Value: v1alpha1.NewProjection("($foo)"),
									}},
								},
								Content: "echo $FOO",
							},
						}},
					},
				}},
			},
		},
	}
	mockTC := func(client client.Client) enginecontext.TestContext {
		registry := mocks.Registry{
			Client: client,
		}
		return enginecontext.MakeContext(clock.RealClock{}, apis.NewBindings(), registry).WithTimeouts(v1alpha1.Timeouts{
			Apply:   &config.Spec.Timeouts.Apply,
			Assert:  &config.Spec.Timeouts.Assert,
			Cleanup: &config.Spec.Timeouts.Cleanup,
			Delete:  &config.Spec.Timeouts.Delete,
			Error:   &config.Spec.Timeouts.Error,
			Exec:    &config.Spec.Timeouts.Exec,
		})
	}
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
		name:   "no tests",
		config: config.Spec,
		tc:     tc(),
		tests:  nil,
		want:   nil,
	}, {
		name:   "test with err",
		config: config.Spec,
		tc:     tc(),
		tests:  []discovery.Test{{Err: errors.New("dummy")}},
		want: &summaryResult{
			failed: 1,
		},
	}, {
		name:   "test with no steps",
		config: config.Spec,
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{{
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{},
			},
		}},
		want: &summaryResult{
			passed: 1,
		},
	}, {
		name:   "skipped test",
		config: config.Spec,
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{{
			Test: &model.Test{
				Spec: v1alpha1.TestSpec{
					Skip: ptr.To(true),
				},
			},
		}},
		want: &summaryResult{
			skipped: 1,
		},
	}, {
		name: "Namesapce exists - success",
		config: model.Configuration{
			Namespace: v1alpha2.NamespaceOptions{
				Name: "default",
			},
		},
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{echoTest},
		want: &summaryResult{
			passed: 1,
		},
	}, {
		name: "Namesapce doesn't exists - success",
		config: model.Configuration{
			Namespace: v1alpha2.NamespaceOptions{
				Name: "chain-saw",
			},
		},
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(v1alpha1.Resource("Namespace"), "chain-saw")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return nil
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{echoTest},
		want: &summaryResult{
			passed: 1,
		},
	}, {
		name: "Namesapce not found - error",
		config: model.Configuration{
			Namespace: v1alpha2.NamespaceOptions{
				Name: "chain-saw",
			},
		},
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewBadRequest("failed to get namespace")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{echoTest},
		want: &summaryResult{
			failed: 1,
		},
	}, {
		name: "Namesapce doesn't exists and can't be created - error",
		config: model.Configuration{
			Namespace: v1alpha2.NamespaceOptions{
				Name: "chain-saw",
			},
		},
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(v1alpha1.Resource("Namespace"), "chain-saw")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return kerrors.NewBadRequest("failed to create namespace")
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{echoTest},
		want: &summaryResult{
			failed: 1,
		},
	}, {
		name: "With namespace compiler",
		config: model.Configuration{
			Namespace: v1alpha2.NamespaceOptions{
				Name:     "chain-saw",
				Compiler: ptr.To(v1alpha2.EngineCEL),
			},
		},
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(v1alpha1.Resource("Namespace"), "chain-saw")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return nil
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{echoTest},
		want: &summaryResult{
			passed: 1,
		},
	}, {
		name: "With scanrio",
		config: model.Configuration{
			Namespace: v1alpha2.NamespaceOptions{
				Name:     "chain-saw",
				Compiler: ptr.To(v1alpha2.EngineCEL),
			},
		},
		tc: func() enginecontext.TestContext {
			client := &fake.FakeClient{
				GetFn: func(ctx context.Context, call int, key ctrlclient.ObjectKey, obj ctrlclient.Object, opts ...ctrlclient.GetOption) error {
					return kerrors.NewNotFound(v1alpha1.Resource("Namespace"), "chain-saw")
				},
				CreateFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.CreateOption) error {
					return nil
				},
				DeleteFn: func(ctx context.Context, call int, obj ctrlclient.Object, opts ...ctrlclient.DeleteOption) error {
					return nil
				},
			}
			return mockTC(client)
		}(),
		tests: []discovery.Test{scenarioTest},
		want: &summaryResult{
			passed: 1,
		},
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
			err := r.Run(context.TODO(), tt.config.Namespace, tt.tc, tt.tests...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil {
				assert.Equal(t, tt.want.failed, tt.tc.Failed())
				assert.Equal(t, tt.want.passed, tt.tc.Passed())
				assert.Equal(t, tt.want.skipped, tt.tc.Skipped())
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
			err := r.run(context.TODO(), tt.m, v1alpha2.NamespaceOptions{}, enginecontext.EmptyContext(clock.RealClock{}), tt.tests...)
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
