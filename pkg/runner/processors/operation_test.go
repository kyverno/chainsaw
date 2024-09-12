package processors

import (
	"context"
	"errors"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine"
	enginecontext "github.com/kyverno/chainsaw/pkg/engine/context"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	mock "github.com/kyverno/chainsaw/pkg/engine/operations/testing"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestOperation_Execute(t *testing.T) {
	tests := []struct {
		name            string
		continueOnError bool
		expectedFail    bool
		operation       operations.Operation
		timeout         time.Duration
	}{{
		name: "operation fails but continues",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ binding.Bindings) (outputs.Outputs, error) {
				return nil, errors.New("operation failed")
			},
		},
		continueOnError: true,
		expectedFail:    true,
		timeout:         1 * time.Second,
	}, {
		name: "operation fails and don't continues",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ binding.Bindings) (outputs.Outputs, error) {
				return nil, errors.New("operation failed")
			},
		},
		continueOnError: false,
		expectedFail:    true,
	}, {
		name: "operation succeeds",
		operation: mock.MockOperation{
			ExecFn: func(_ context.Context, _ binding.Bindings) (outputs.Outputs, error) {
				return nil, nil
			},
		},
		expectedFail: false,
		timeout:      1 * time.Second,
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			localTC := tc
			op := newOperation(
				OperationInfo{},
				localTC.continueOnError,
				func(ctx context.Context, tc engine.Context) (operations.Operation, *time.Duration, engine.Context, error) {
					return localTC.operation, &localTC.timeout, tc, nil
				},
			)
			nt := testing.MockT{}
			ctx := testing.IntoContext(context.Background(), &nt)
			tcontext := enginecontext.EmptyContext()
			op.execute(ctx, tcontext)
			if localTC.expectedFail {
				assert.True(t, nt.FailedVar, "expected an error but got none")
			} else {
				assert.False(t, nt.FailedVar, "expected no error but got one")
			}
		})
	}
}
