package sleep

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/logging"
)

type operation struct {
	duration time.Duration
}

func New(duration time.Duration) operations.Operation {
	return &operation{
		duration: duration,
	}
}

func (o *operation) Exec(ctx context.Context, _ apis.Bindings) (_ outputs.Outputs, _err error) {
	defer func() {
		internal.LogEnd(ctx, logging.Sleep, nil, _err)
	}()
	internal.LogStart(ctx, logging.Sleep, nil)
	return nil, o.execute()
}

func (o *operation) execute() error {
	time.Sleep(o.duration)
	return nil
}
