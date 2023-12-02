package processors

import (
	"context"
	"errors"
	"time"

	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	mock "github.com/kyverno/chainsaw/pkg/runner/operations/testing"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestOperation_Execute(t *testing.T) {
	tests := []struct {
		name            string
		continueOnError bool
		expectedFail    bool
		operation       operations.Operation
		operationReport *report.OperationReport
	}{
		{
			name: "operation fails but continues",
			operation: mock.MockOperation{
				ExecFn: func(ctx context.Context) error {
					return errors.New("operation failed")
				},
			},
			continueOnError: true,
			expectedFail:    true,
			operationReport: report.NewOperation("FakeOperation", report.OperationTypeCreate),
		},
		// {
		// 	name: "operation fails don't continues",
		// 	operation: mock.MockOperation{
		// 		ExecFn: func(ctx context.Context) error {
		// 			return errors.New("operation failed")
		// 		},
		// 	},
		// 	continueOnError: false,
		// 	expectedFail:    true,
		// 	operationReport: report.NewOperation("FakeOperation", report.OperationTypeCreate),
		// },
		{
			name: "operation succeeds",
			operation: mock.MockOperation{
				ExecFn: func(ctx context.Context) error {
					return nil
				},
			},
			expectedFail:    false,
			operationReport: report.NewOperation("FakeOperation", report.OperationTypeCreate),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			op := operation{
				continueOnError: tc.continueOnError,
				timeout:         time.Duration(1),
				operation:       tc.operation,
				operationReport: tc.operationReport,
			}
			nt := testing.T{}
			ctx := testing.IntoContext(context.Background(), &nt)
			op.execute(ctx)

			if tc.expectedFail {
				assert.True(t, nt.Failed(), "expected an error but got none")
			} else {
				assert.False(t, nt.Failed(), "expected no error but got one")
			}
		})
	}
}
