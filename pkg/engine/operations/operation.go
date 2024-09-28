package operations

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
)

type Operation interface {
	Exec(context.Context, apis.Bindings) (outputs.Outputs, error)
}
