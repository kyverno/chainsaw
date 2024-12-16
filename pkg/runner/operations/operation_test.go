package operations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	mock "github.com/kyverno/chainsaw/pkg/engine/operations/testing"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/clock"
)

func TestOperation_Execute(t *testing.T) {
	tests := []struct {
		name         string
		expectedFail bool
		operation    operations.Operation
		timeout      time.Duration
	}{{
		name: "operation fails but continues",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ apis.Bindings) (outputs.Outputs, error) {
				return nil, errors.New("operation failed")
			},
		},
		expectedFail: true,
		timeout:      1 * time.Second,
	}, {
		name: "operation fails and don't continues",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ apis.Bindings) (outputs.Outputs, error) {
				return nil, errors.New("operation failed")
			},
		},
		expectedFail: true,
	}, {
		name: "operation succeeds",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ apis.Bindings) (outputs.Outputs, error) {
				return nil, nil
			},
		},
		timeout: 1 * time.Second,
	}, {
		name: "operation succeeds",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ apis.Bindings) (outputs.Outputs, error) {
				return nil, nil
			},
		},
		timeout: 1 * time.Second,
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			op := newOperation(
				func(ctx context.Context, _tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
					return tc.operation, &tc.timeout, _tc, nil
				},
			)
			tcontext := enginecontext.EmptyContext(clock.RealClock{})
			ctx := context.Background()
			_, err := op.Execute(ctx, tcontext, &model.StepReport{})
			if tc.expectedFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	{
		t.Run("factory failure", func(t *testing.T) {
			op := newOperation(
				func(ctx context.Context, tc enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error) {
					return nil, nil, tc, errors.New("dummy")
				},
			)
			tcontext := enginecontext.EmptyContext(clock.RealClock{})
			ctx := context.Background()
			_, err := op.Execute(ctx, tcontext, &model.StepReport{})
			assert.Error(t, err)
		})
	}
}
