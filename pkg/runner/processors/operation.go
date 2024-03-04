package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
)

type operation struct {
	continueOnError bool
	timeout         *time.Duration
	operation       func() (operations.Operation, error)
	operationReport *report.OperationReport
	bindings        binding.Bindings
	variables       []v1alpha1.Binding
}

func newOperation(
	continueOnError bool,
	timeout *time.Duration,
	op operations.Operation,
	operationReport *report.OperationReport,
	bindings binding.Bindings,
	variables ...v1alpha1.Binding,
) operation {
	return operation{
		continueOnError: continueOnError,
		timeout:         timeout,
		operation: func() (operations.Operation, error) {
			return op, nil
		},
		operationReport: operationReport,
		bindings:        bindings,
		variables:       variables,
	}
}

func newLazyOperation(
	continueOnError bool,
	timeout *time.Duration,
	op func() (operations.Operation, error),
	operationReport *report.OperationReport,
	bindings binding.Bindings,
	variables ...v1alpha1.Binding,
) operation {
	return operation{
		continueOnError: continueOnError,
		timeout:         timeout,
		operation:       op,
		operationReport: operationReport,
		bindings:        bindings,
		variables:       variables,
	}
}

func (o operation) execute(ctx context.Context) {
	if o.timeout != nil {
		toCtx, cancel := context.WithTimeout(ctx, *o.timeout)
		ctx = toCtx
		defer cancel()
	}
	handleError := func(err error) {
		t := testing.FromContext(ctx)
		if err != nil {
			logging.Log(ctx, logging.Internal, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		}
		if o.continueOnError {
			t.Fail()
		} else {
			t.FailNow()
		}
	}
	operation, err := o.operation()
	if err != nil {
		handleError(err)
	} else if bindings, err := registerBindings(ctx, o.bindings, o.variables...); err != nil {
		handleError(err)
	} else {
		err := operation.Exec(ctx, bindings)
		if o.operationReport != nil {
			o.operationReport.MarkOperationEnd(err)
		}
		if err != nil {
			handleError(err)
		}
	}
}
