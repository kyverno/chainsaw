package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
)

type Operation interface {
	Exec(ctx context.Context) error
	Cleanup()
}

type baseOperation struct {
	client client.Client
}

func execOperation(ctx context.Context, operation Operation) error {
	err := operation.Exec(ctx)
	operation.Cleanup()
	return err
}
