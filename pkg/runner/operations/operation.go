package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/engine/outputs"
	enginecontext "github.com/kyverno/chainsaw/pkg/runner/context"
)

type Operation interface {
	Execute(context.Context, enginecontext.TestContext) (outputs.Outputs, error)
}
