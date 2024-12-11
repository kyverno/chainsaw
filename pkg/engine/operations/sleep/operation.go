package sleep

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/operations"
	"github.com/kyverno/chainsaw/pkg/engine/operations/internal"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	"github.com/kyverno/chainsaw/pkg/logging"
)

type operation struct {
	duration v1alpha1.Sleep
}

func New(duration v1alpha1.Sleep) operations.Operation {
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
	time.Sleep(o.duration.Duration.Duration)
	return nil
}
