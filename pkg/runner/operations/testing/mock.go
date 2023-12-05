package testing

import "context"

type MockOperation struct {
	ExecFn func(ctx context.Context) error
}

func (m MockOperation) Exec(ctx context.Context) error {
	return m.ExecFn(ctx)
}
