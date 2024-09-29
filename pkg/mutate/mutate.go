package mutate

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Mutate(ctx context.Context, compilers compilers.Compilers, path *field.Path, mutation Mutation, value any, bindings apis.Bindings) (any, error) {
	return mutation.mutate(ctx, compilers, path, value, bindings)
}
