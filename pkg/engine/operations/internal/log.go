package internal

import (
	"context"
	"fmt"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/logging"
	"github.com/kyverno/pkg/ext/output/color"
)

func LogStart(ctx context.Context, op logging.Operation, obj client.Object, args ...fmt.Stringer) {
	logging.Log(ctx, op, logging.RunStatus, obj, color.BoldFgCyan, args...)
}

func LogEnd(ctx context.Context, op logging.Operation, obj client.Object, err error) {
	if err != nil {
		logging.Log(ctx, op, logging.ErrorStatus, obj, color.BoldRed, logging.ErrSection(err))
	} else {
		logging.Log(ctx, op, logging.DoneStatus, obj, color.BoldGreen)
	}
}
