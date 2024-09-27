package testing

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
)

type MockOperation struct {
	ExecFn func(context.Context, apis.Bindings) (outputs.Outputs, error)
}

func (m MockOperation) Exec(ctx context.Context, bindings apis.Bindings) (outputs.Outputs, error) {
	return m.ExecFn(ctx, bindings)
}
