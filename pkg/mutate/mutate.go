package mutate

import (
	"context"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/kyverno-json/pkg/engine/template"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Mutate(ctx context.Context, path *field.Path, mutation Mutation, value any, bindings binding.Bindings, opts ...template.Option) (any, error) {
	return mutation.mutate(ctx, path, value, bindings, opts...)
}
