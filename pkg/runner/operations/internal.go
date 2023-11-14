package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/runner/logging"
	"github.com/kyverno/kyverno/ext/output/color"
)

type Operation interface {
	Exec(ctx context.Context) error
	Name() string
}

type BaseOperation struct {
	client client.Client
}

func execOperation(ctx context.Context, operation Operation) error {
	logger := logging.FromContext(ctx)
	logger.Log("Starting operation %s ", color.BoldFgCyan, operation.Name())
	err := operation.Exec(ctx)
	if err != nil {
		logger.Log("Operation %s failed with error %s", color.BoldRed, operation.Name(), err)
	} else {
		logger.Log("Operation %s completed successfully", color.BoldGreen, operation.Name())
	}

	return err
}
