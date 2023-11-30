package testing

import "context"

type MockOperation struct {
	numCalls int
	ExecFn   func(ctx context.Context) error
}

func (m MockOperation) Exec(ctx context.Context) error {
	return m.ExecFn(ctx)
}

func (m *MockOperation) NumCalls() int {
	return m.numCalls
}
