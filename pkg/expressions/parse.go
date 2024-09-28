package expressions

import (
	"context"
	"regexp"

	"github.com/kyverno/kyverno-json/pkg/core/expression"
)

var (
	escapeRegex = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex = regexp.MustCompile(`^\((?:(\w+);)?(.+)\)$`)
)

type Expression struct {
	Statement string
	Engine    string
}

func Parse(ctx context.Context, value string) *Expression {
	return parseExpressionRegex(ctx, value)
}

func parseExpressionRegex(_ context.Context, in string) *Expression {
	out := &Expression{}
	// 1. match escape, if there's no escaping then match engine
	if match := escapeRegex.FindStringSubmatch(in); match != nil {
		in = match[1]
	} else {
		if match := engineRegex.FindStringSubmatch(in); match != nil {
			out.Engine = match[1]
			// account for default engine
			if out.Engine == "" {
				out.Engine = expression.CompilerDefault
			}
			in = match[2]
		}
	}
	// parse statement
	out.Statement = in
	if out.Statement == "" {
		return nil
	}
	return out
}
