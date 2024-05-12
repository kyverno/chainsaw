package processors

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/report"
	apibindings "github.com/kyverno/chainsaw/pkg/runner/bindings"
	"github.com/kyverno/chainsaw/pkg/runner/clusters"
	"github.com/kyverno/chainsaw/pkg/runner/failer"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/pkg/ext/output/color"
	"k8s.io/client-go/rest"
)

type operation struct {
	cluster         clusters.Cluster
	info            OperationInfo
	continueOnError bool
	timeout         *time.Duration
	operation       func(context.Context, binding.Bindings) (operations.Operation, error)
	report          *report.OperationReport
	variables       []v1alpha1.Binding
}

func newLazyOperation(
	cluster clusters.Cluster,
	info OperationInfo,
	continueOnError bool,
	timeout *time.Duration,
	op func(context.Context, binding.Bindings) (operations.Operation, error),
	report *report.OperationReport,
	variables ...v1alpha1.Binding,
) operation {
	return operation{
		cluster:         cluster,
		info:            info,
		continueOnError: continueOnError,
		timeout:         timeout,
		operation:       op,
		report:          report,
		variables:       variables,
	}
}

func (o operation) execute(ctx context.Context, bindings binding.Bindings) operations.Outputs {
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
	operation, err := o.operation(ctx, bindings)
	var config *rest.Config
	var client client.Client
	if o.cluster != nil {
		config = o.cluster.Config()
		client = o.cluster.Client()
	}
	bindings = apibindings.RegisterClusterBindings(ctx, bindings, config, client)
	if err != nil {
		handleError(err)
	} else if bindings, err := apibindings.RegisterBindings(ctx, bindings, o.variables...); err != nil {
		handleError(err)
	} else {
		outputs, err := operation.Exec(ctx, apibindings.RegisterNamedBinding(ctx, bindings, "operation", o.info))
		// TODO
		// if o.operationReport != nil {
		// 	o.operationReport.MarkOperationEnd(err)
		// }
		if err != nil {
			handleError(nil)
		}
		return outputs
	}
	return nil
}
