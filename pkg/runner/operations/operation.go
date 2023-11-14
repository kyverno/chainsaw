package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
)

type Operation interface {
	Exec(ctx context.Context) error
	Name() string
}

type BaseOperation struct {
	client client.Client
}

func execOperation(ctx context.Context, operation Operation) error {
	return operation.Exec(ctx)
}
