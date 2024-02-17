package sleep

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/runner/operations/internal"
)

type operation struct {
	duration v1alpha1.Sleep
}

func New(duration v1alpha1.Sleep) operations.Operation {
	return &operation{
		duration: duration,
	}
}

func (o *operation) Exec(ctx context.Context) (err error) {
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Sleep, err)
	}()
	internal.LogStart(logger, logging.Sleep)
	return o.execute(ctx)
}

func (o *operation) execute(ctx context.Context) error {
	time.Sleep(o.duration.Duration.Duration)
	return nil
}
