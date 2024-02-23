package operations

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
)

type Operation interface {
	Exec(context.Context, binding.Bindings) error
}
