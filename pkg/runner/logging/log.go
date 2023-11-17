package logging

import (
	"context"
	"fmt"

	"github.com/kyverno/kyverno/ext/output/color"
)

func Log(ctx context.Context, operation Operation, status Status, color *color.Color, args ...fmt.Stringer) {
	logger := FromContext(ctx)
	if logger != nil {
		logger.Log(operation, status, color, args...)
	}
}
