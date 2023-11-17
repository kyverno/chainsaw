package processors

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/runner/operations"
	"github.com/kyverno/chainsaw/pkg/testing"
)

type operation struct {
	continueOnError bool
	timeout         time.Duration
	operation       operations.Operation
}

func (o operation) execute(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()
	if err := o.operation.Exec(ctx); err != nil {
		t := testing.FromContext(ctx)
		if o.continueOnError {
			t.Fail()
		} else {
			t.FailNow()
		}
	}
}
