package operations

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	opsleep "github.com/kyverno/chainsaw/pkg/engine/operations/sleep"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

type sleepAction struct {
	duration time.Duration
}

func (o sleepAction) Execute(ctx context.Context, tc enginecontext.TestContext) (outputs.Outputs, error) {
	op := opsleep.New(o.duration)
	return op.Exec(ctx, tc.Bindings())
}

func sleepOperation(op v1alpha1.Sleep) Operation {
	return sleepAction{
		duration: op.Duration.Duration,
	}
}
