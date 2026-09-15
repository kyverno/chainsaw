package expressions

import (
	"context"
	"regexp"
	"strings"

	"github.com/kyverno/kyverno-json/pkg/core/expression"
)

var (
	escapeRegex    = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex    = regexp.MustCompile(`^\((?:(\w+);)?(.+)\)$`)
	celBindingName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*`)
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
	// for CEL statements, rewrite the `$name` shorthand into an explicit
	// `bindings.resolve('name')` call so users don't have to spell that
	// out themselves
	if out.Engine == expression.CompilerCEL {
		in = celBindingShorthand(in)
	}
	// parse statement
	out.Statement = in
	if out.Statement == "" {
		return nil
	}
	return out
}

// celBindingShorthand rewrites `$name` occurrences in a CEL statement into
// `bindings.resolve('name')`, skipping anything inside a single- or
// double-quoted string literal so a literal `$` in a string (e.g. "cost:
// $5") is left untouched.
func celBindingShorthand(in string) string {
	var out strings.Builder
	var quote rune
	runes := []rune(in)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if quote != 0 {
			out.WriteRune(c)
			if c == '\\' && i+1 < len(runes) {
				// preserve the escaped character verbatim so it can't
				// prematurely close the string literal
				i++
				out.WriteRune(runes[i])
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			out.WriteRune(c)
			continue
		}
		if c == '$' {
			if name := celBindingName.FindString(string(runes[i+1:])); name != "" {
				out.WriteString("bindings.resolve('")
				out.WriteString(name)
				out.WriteString("')")
				i += len(name)
				continue
			}
		}
		out.WriteRune(c)
	}
	return out.String()
}
