package testing

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
)

type MockOperation struct {
	ExecFn func(context.Context, binding.Bindings) (outputs.Outputs, error)
}

func (m MockOperation) Exec(ctx context.Context, bindings binding.Bindings) (outputs.Outputs, error) {
	return m.ExecFn(ctx, bindings)
}
