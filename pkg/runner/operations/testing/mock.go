package testing

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/runner/operations"
)

type MockOperation struct {
	ExecFn func(context.Context, binding.Bindings) (operations.Outputs, error)
}

func (m MockOperation) Exec(ctx context.Context, bindings binding.Bindings) (operations.Outputs, error) {
	return m.ExecFn(ctx, bindings)
}
