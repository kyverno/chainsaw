package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/report"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/testing"
	"github.com/kyverno/kyverno/ext/output/color"
	"k8s.io/client-go/rest"
)

type operation struct {
	continueOnError bool
	timeout         *time.Duration
	operation       func() (operations.Operation, error)
	operationReport *report.OperationReport
	config          *rest.Config
	client          client.Client
	variables       []v1alpha1.Binding
}

func newOperation(
	continueOnError bool,
	timeout *time.Duration,
	op operations.Operation,
	operationReport *report.OperationReport,
	config *rest.Config,
	client client.Client,
	variables ...v1alpha1.Binding,
) operation {
	return newLazyOperation(
		continueOnError,
		timeout,
		func() (operations.Operation, error) {
			return op, nil
		},
		operationReport,
		config,
		client,
		variables...,
	)
}

func newLazyOperation(
	continueOnError bool,
	timeout *time.Duration,
	op func() (operations.Operation, error),
	operationReport *report.OperationReport,
	config *rest.Config,
	client client.Client,
	variables ...v1alpha1.Binding,
) operation {
	return operation{
		continueOnError: continueOnError,
		timeout:         timeout,
		operation:       op,
		operationReport: operationReport,
		client:          client,
		config:          config,
		variables:       variables,
	}
}

func (o operation) execute(ctx context.Context, bindings binding.Bindings) operations.Outputs {
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
	} else if bindings, err := registerBindings(ctx, bindings, o.config, o.client, o.variables...); err != nil {
		handleError(err)
	} else {
		outputs, err := operation.Exec(ctx, bindings)
		if o.operationReport != nil {
			o.operationReport.MarkOperationEnd(err)
		}
		if err != nil {
			handleError(nil)
		}
		return outputs
	}
	return nil
}
