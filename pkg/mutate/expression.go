package mutate

import (
	"context"
	"reflect"

	"github.com/kyverno/chainsaw/pkg/expressions"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

func parseExpression(ctx context.Context, value any) *expressions.Expression {
	if reflectutils.GetKind(value) != reflect.String {
		return nil
	}
	return expressions.Parse(ctx, reflect.ValueOf(value).String())
}
