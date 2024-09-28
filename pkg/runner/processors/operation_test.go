package processors

import (
	"context"
	"errors"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine"
	enginecontext "github.com/kyverno/chainsaw/pkg/engine/context"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	mock "github.com/kyverno/chainsaw/pkg/engine/operations/testing"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
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
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			localTC := tc
			op := newOperation(
				OperationInfo{},
				model.OperationTypeApply,
				func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
					return localTC.operation, &localTC.timeout, tc, nil
				},
			)
			tcontext := enginecontext.EmptyContext()
			ctx := context.Background()
			_, err := op.execute(ctx, tcontext, &model.StepReport{})
			if localTC.expectedFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
