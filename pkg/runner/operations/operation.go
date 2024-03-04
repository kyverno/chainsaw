package operations

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type Outputs = map[string]binding.Binding

type Operation interface {
	Exec(context.Context, binding.Bindings) (Outputs, error)
}
