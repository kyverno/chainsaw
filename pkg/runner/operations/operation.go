package operations

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/logging"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/pkg/ext/output/color"
)

type Operation interface {
	Execute(context.Context, enginecontext.TestContext) (outputs.Outputs, error)
}

type operationFactory = func(context.Context, enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error)

type operation struct {
	operation operationFactory
}

func newOperation(
	op operationFactory,
) Operation {
	return operation{
		operation: op,
	}
}

func (o operation) Execute(ctx context.Context, tc enginecontext.TestContext) (_ outputs.Outputs, err error) {
	if operation, timeout, tc, err := o.operation(ctx, tc); err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, nil, color.BoldRed, logging.ErrSection(err))
		return nil, err
	} else {
		if timeout != nil && *timeout != 0 {
			toCtx, cancel := context.WithTimeout(ctx, *timeout)
			ctx = toCtx
			defer cancel()
		}
		return operation.Exec(ctx, tc.Bindings())
	}
}
