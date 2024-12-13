package runner

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
	"github.com/kyverno/pkg/ext/output/color"
)

type operationFactory = func(context.Context, enginecontext.TestContext) (operations.Operation, *time.Duration, enginecontext.TestContext, error)

type operation struct {
	info      OperationInfo
	opType    model.OperationType
	operation operationFactory
}

func newOperation(
	info OperationInfo,
	opType model.OperationType,
	op operationFactory,
) operation {
	return operation{
		info:      info,
		opType:    opType,
		operation: op,
	}
}

func (o operation) execute(ctx context.Context, tc enginecontext.TestContext, stepReport *model.StepReport) (_ outputs.Outputs, err error) {
	report := &model.OperationReport{
		Type:      o.opType,
		StartTime: time.Now(),
	}
	defer func() {
		report.EndTime = time.Now()
		report.Err = err
		stepReport.Add(report)
	}()
	if operation, timeout, tc, err := o.operation(ctx, tc.WithBinding("operation", o.info)); err != nil {
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
