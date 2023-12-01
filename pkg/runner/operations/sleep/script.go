package script

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/kyverno/ext/output/color"
)

type operation struct {
	sleep v1alpha1.Sleep
}

func New(sleep v1alpha1.Sleep) operations.Operation {
	return &operation{
		sleep: sleep,
	}
}

func (o *operation) Exec(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	defer func() {
		logger.Log(logging.Sleep, logging.DoneStatus, color.BoldGreen)
	}()
	logger.Log(logging.Sleep, logging.RunStatus, color.BoldFgCyan)
	time.Sleep(o.sleep.Duration.Duration)
	return nil
}
