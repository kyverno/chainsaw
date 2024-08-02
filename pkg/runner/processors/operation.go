package processors

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/engine"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/pkg/ext/output/color"
)

type operationFactory = func(context.Context, engine.Context) (operations.Operation, *time.Duration, engine.Context, error)

type operation struct {
	info            OperationInfo
	continueOnError bool
	operation       operationFactory
	report          *report.OperationReport
}

func newOperation(
	info OperationInfo,
	continueOnError bool,
	op operationFactory,
	report *report.OperationReport,
) operation {
	return operation{
		info:            info,
		continueOnError: continueOnError,
		operation:       op,
		report:          report,
	}
}

func (o operation) execute(ctx context.Context, tc engine.Context) outputs.Outputs {
	if o.report != nil {
		o.report.SetStartTime(time.Now())
		defer func() {
			o.report.SetEndTime(time.Now())
		}()
	}
	handleError := func(err error) {
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		}
		if o.continueOnError {
			failer.Fail(ctx)
		} else {
			failer.FailNow(ctx)
		}
	}
	tc = tc.WithBinding(ctx, "operation", o.info)
	operation, timeout, tc, err := o.operation(ctx, tc)
	if err != nil {
		handleError(err)
	} else {
		if timeout != nil {
			toCtx, cancel := context.WithTimeout(ctx, *timeout)
			ctx = toCtx
			defer cancel()
		}
		outputs, err := operation.Exec(ctx, tc.Bindings())
		// TODO
		// if o.operationReport != nil {
		// 	o.operationReport.MarkOperationEnd(err)
		// }
		if err != nil {
			// we pass nil in the err argument so that it is not logged in the output
			handleError(nil)
		}
		return outputs
	}
	return nil
}
