package logging

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/kyverno/chainsaw/pkg/client"
)

func Log(ctx context.Context, operation Operation, status Status, obj client.Object, color *color.Color, args ...fmt.Stringer) {
	if logger := getLogger(ctx); logger != nil {
		logger.Log(ctx, operation, status, obj, color, args...)
	}
}
