package operations

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type Outputs = map[string]any

type Operation interface {
	Exec(context.Context, binding.Bindings) (Outputs, error)
}
