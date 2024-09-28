package mutate

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Mutate(ctx context.Context, path *field.Path, mutation Mutation, value any, bindings apis.Bindings, opts ...jp.Option) (any, error) {
	return mutation.mutate(ctx, path, value, bindings, opts...)
}
