package internal

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetLogger(ctx context.Context, obj client.Object) logging.Logger {
	return logging.FromContext(ctx).WithResource(obj)
}

func LogStart(logger logging.Logger, op logging.Operation) {
	logger.Log(op, logging.RunStatus, color.BoldFgCyan)
}

func LogEnd(logger logging.Logger, op logging.Operation, err error) {
	if err != nil {
		logger.Log(op, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
	} else {
		logger.Log(op, logging.DoneStatus, color.BoldGreen)
	}
}
