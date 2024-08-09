package expressions

import (
	"context"
	"reflect"
	"regexp"
)

var (
	escapeRegex = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex = regexp.MustCompile(`^\((?:(\w+):)?(.+)\)$`)
)

type Expression struct {
	Statement string
	Engine    string
}

func Parse(ctx context.Context, value string) *Expression {
	return parseExpressionRegex(ctx, reflect.ValueOf(value).String())
}

func parseExpressionRegex(_ context.Context, in string) *Expression {
	expression := &Expression{}
	// 1. match escape, if there's no escaping then match engine
	if match := escapeRegex.FindStringSubmatch(in); match != nil {
		in = match[1]
	} else {
		if match := engineRegex.FindStringSubmatch(in); match != nil {
			expression.Engine = match[1]
			// account for default engine
			if expression.Engine == "" {
				expression.Engine = "jp"
			}
			in = match[2]
		}
	}
	// parse statement
	expression.Statement = in
	if expression.Statement == "" {
		return nil
	}
	return expression
}
