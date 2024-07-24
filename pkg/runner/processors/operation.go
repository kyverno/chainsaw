package processors

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/chainsaw/pkg/model"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/pkg/ext/output/color"
)

type operation struct {
	info            OperationInfo
	continueOnError bool
	timeout         *time.Duration
	operation       func(context.Context, model.TestContext) (operations.Operation, model.TestContext, error)
	report          *report.OperationReport
}

func newOperation(
	info OperationInfo,
	continueOnError bool,
	timeout *time.Duration,
	op func(context.Context, model.TestContext) (operations.Operation, model.TestContext, error),
	report *report.OperationReport,
) operation {
	return operation{
		info:            info,
		continueOnError: continueOnError,
		timeout:         timeout,
		operation:       op,
		report:          report,
	}
}

func (o operation) execute(ctx context.Context, tc model.TestContext) operations.Outputs {
	if o.report != nil {
		o.report.SetStartTime(time.Now())
		defer func() {
			o.report.SetEndTime(time.Now())
		}()
	}
	if o.timeout != nil {
		toCtx, cancel := context.WithTimeout(ctx, *o.timeout)
		ctx = toCtx
		defer cancel()
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
	operation, tc, err := o.operation(ctx, tc)
	if err != nil {
		handleError(err)
	} else {
		outputs, err := operation.Exec(ctx, apibindings.RegisterNamedBinding(ctx, tc.Bindings(), "operation", o.info))
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
