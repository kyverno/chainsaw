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

	return operation.Exec(ctx)
}
