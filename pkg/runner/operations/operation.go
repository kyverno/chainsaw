package operations

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/engine/outputs"
)

type Operation interface {
	Exec(context.Context, binding.Bindings) (outputs.Outputs, error)
}
