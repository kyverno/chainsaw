package mutate

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/kyverno-json/pkg/core/compilers/jp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Mutation interface {
	mutate(context.Context, *field.Path, any, apis.Bindings, ...jp.Option) (any, error)
}
