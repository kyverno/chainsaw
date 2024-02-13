package mutate

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Mutation interface {
	mutate(context.Context, *field.Path, any, binding.Bindings, ...template.Option) (any, error)
}
