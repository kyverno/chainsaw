package logging

import (
	"context"

	"github.com/kyverno/kyverno/ext/output/color"
)

func Log(ctx context.Context, operation Operation, color *color.Color, args ...interface{}) {
	logger := FromContext(ctx)
	if logger != nil {
		logger.Log(operation, color, args...)
	}
}
