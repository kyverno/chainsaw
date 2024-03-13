package mutate

import (
	"context"
	"reflect"
	"regexp"

	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

var (
	escapeRegex = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex = regexp.MustCompile(`^\((?:(\w+):)?(.+)\)$`)
)

type expression struct {
	statement string
	engine    string
}

func parseExpressionRegex(_ context.Context, in string) *expression {
	expression := &expression{}
	// 1. match escape, if there's no escaping then match engine
	if match := escapeRegex.FindStringSubmatch(in); match != nil {
		in = match[1]
	} else {
		if match := engineRegex.FindStringSubmatch(in); match != nil {
			expression.engine = match[1]
			// account for default engine
			if expression.engine == "" {
				expression.engine = "jp"
			}
			in = match[2]
		}
	}
	// parse statement
	expression.statement = in
	if expression.statement == "" {
		return nil
	}
	return expression
}

func parseExpression(ctx context.Context, value any) *expression {
	if reflectutils.GetKind(value) != reflect.String {
		return nil
	}
	return parseExpressionRegex(ctx, reflect.ValueOf(value).String())
}
