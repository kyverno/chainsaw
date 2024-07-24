package sleep

import (
	"context"
	"time"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/logging"
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

func (o *operation) Exec(ctx context.Context, _ binding.Bindings) (_ operations.Outputs, _err error) {
	logger := internal.GetLogger(ctx, nil)
	defer func() {
		internal.LogEnd(logger, logging.Sleep, _err)
	}()
	internal.LogStart(logger, logging.Sleep)
	return nil, o.execute()
}

func (o *operation) execute() error {
	time.Sleep(o.duration.Duration.Duration)
	return nil
}
