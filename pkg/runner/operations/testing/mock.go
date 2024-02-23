package testing

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type MockOperation struct {
	ExecFn func(context.Context, binding.Bindings) error
}

func (m MockOperation) Exec(ctx context.Context, bindings binding.Bindings) error {
	return m.ExecFn(ctx, bindings)
}
