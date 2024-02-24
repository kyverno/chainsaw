package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/testing"
)

type operation struct {
	continueOnError bool
	timeout         *time.Duration
	operation       operations.Operation
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
	handleError := func() {
		t := testing.FromContext(ctx)
		if o.continueOnError {
			t.Fail()
		} else {
			t.FailNow()
		}
	}
	if bindings, err := registerBindings(ctx, o.bindings, o.variables...); err != nil {
		handleError()
	} else {
		err := o.operation.Exec(ctx, bindings)
		if o.operationReport != nil {
			o.operationReport.MarkOperationEnd(err)
		}
		if err != nil {
			handleError()
		}
	}
}
