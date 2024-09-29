package mutate

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Mutate(ctx context.Context, path *field.Path, mutation Mutation, value any, bindings apis.Bindings, compilers compilers.Compilers) (any, error) {
	return mutation.mutate(ctx, path, value, bindings, compilers)
}
