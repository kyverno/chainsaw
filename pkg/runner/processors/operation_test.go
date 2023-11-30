package processors

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/runner/operations"
	mock "github.com/kyverno/chainsaw/pkg/runner/operations/testing"
	"github.com/kyverno/chainsaw/pkg/testing"
)

func TestOperation_Execute(t *testing.T) {
	tests := []struct {
		name            string
		continueOnError bool
		expectedFail    bool
		operation       operations.Operation
	}{
		// {
		// 	name: "operation fails",
		// 	operation: mock.MockOperation{
		// 		ExecFn: func(ctx context.Context) error {
		// 			return errors.New("operation failed")
		// 		},
		// 	},
		// 	expectedFail: true,
		// },
		{
			name: "operation succeeds",
			operation: mock.MockOperation{
				ExecFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedFail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			op := operation{
				continueOnError: tc.continueOnError,
				timeout:         1 * time.Second,
				operation:       tc.operation,
				operationReport: nil,
			}
			op.execute(context.Background())
		})
	}
}
