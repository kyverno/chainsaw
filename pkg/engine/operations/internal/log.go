package internal

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
	"github.com/kyverno/pkg/ext/output/color"
)

func GetLogger(ctx context.Context, obj client.Object) logging.Logger {
	logger := logging.FromContext(ctx)
	if logger == nil {
		return logger
	}
	if obj != nil {
		if obj.GetObjectKind().GroupVersionKind().Kind == "" {
			return logger
		}
	}
	return logger.WithResource(obj)
}

func LogStart(logger logging.Logger, op logging.Operation, args ...fmt.Stringer) {
	if logger != nil {
		logger.Log(op, logging.RunStatus, color.BoldFgCyan, args...)
	}
}

func LogEnd(logger logging.Logger, op logging.Operation, err error) {
	if logger != nil {
		if err != nil {
			logger.Log(op, logging.ErrorStatus, color.BoldRed, logging.ErrSection(err))
		} else {
			logger.Log(op, logging.DoneStatus, color.BoldGreen)
		}
	}
}
