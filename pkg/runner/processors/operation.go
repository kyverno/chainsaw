package processors

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/pkg/ext/output/color"
)

type operationFactory = func(context.Context, engine.Context) (operations.Operation, *time.Duration, engine.Context, error)

type operation struct {
	info      OperationInfo
	operation operationFactory
}

func newOperation(
	info OperationInfo,
	op operationFactory,
) operation {
	return operation{
		info:      info,
		operation: op,
	}
}

func (o operation) execute(ctx context.Context, tc engine.Context, stepReport *model.StepReport) (outputs.Outputs, error) {
	report := model.OperationReport{
		StartTime: time.Now(),
	}
	defer func() {
		report.EndTime = time.Now()
		stepReport.Add(report)
	}()
	tc = tc.WithBinding(ctx, "operation", o.info)
	operation, timeout, tc, err := o.operation(ctx, tc)
	if err != nil {
		logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		return nil, err
	} else {
		if timeout != nil {
			toCtx, cancel := context.WithTimeout(ctx, *timeout)
			ctx = toCtx
			defer cancel()
		}
		return operation.Exec(ctx, tc.Bindings())
	}
}
