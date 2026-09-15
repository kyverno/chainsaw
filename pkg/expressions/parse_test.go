package expressions

import (
	"context"
	"testing"

	"github.com/kyverno/kyverno-json/pkg/core/expression"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantNil  bool
		wantStmt string
		wantEng  string
	}{{
		name:    "empty",
		in:      "",
		wantNil: true,
	}, {
		name:     "escaped",
		in:       `\(foo)\`,
		wantStmt: "(foo)",
		wantEng:  "",
	}, {
		name:     "default engine, no shorthand rewrite",
		in:       "(foo)",
		wantStmt: "foo",
		wantEng:  expression.CompilerDefault,
	}, {
		name:     "explicit engine, no shorthand rewrite outside cel",
		in:       "(jp;$foo)",
		wantStmt: "$foo",
		wantEng:  "jp",
	}, {
		name:     "cel, single binding shorthand",
		in:       "(cel;$foo == 'bar')",
		wantStmt: "bindings.resolve('foo') == 'bar'",
		wantEng:  expression.CompilerCEL,
	}, {
		name:     "cel, multiple binding shorthands",
		in:       "(cel;$outputs_result.Namespaces.exists(ns, ns.ID==$namespace_id))",
		wantStmt: "bindings.resolve('outputs_result').Namespaces.exists(ns, ns.ID==bindings.resolve('namespace_id'))",
		wantEng:  expression.CompilerCEL,
	}, {
		name:     "cel, dollar inside single-quoted string is untouched",
		in:       "(cel;'cost: $5' == $price)",
		wantStmt: "'cost: $5' == bindings.resolve('price')",
		wantEng:  expression.CompilerCEL,
	}, {
		name:     "cel, dollar inside double-quoted string is untouched",
		in:       `(cel;"cost: $5" == $price)`,
		wantStmt: `"cost: $5" == bindings.resolve('price')`,
		wantEng:  expression.CompilerCEL,
	}, {
		name:     "cel, escaped quote inside string doesn't close it early",
		in:       `(cel;"it\'s $5" == $price)`,
		wantStmt: `"it\'s $5" == bindings.resolve('price')`,
		wantEng:  expression.CompilerCEL,
	}, {
		name:     "cel, already-explicit resolve call is left as-is",
		in:       "(cel;bindings.resolve('foo') == 'bar')",
		wantStmt: "bindings.resolve('foo') == 'bar'",
		wantEng:  expression.CompilerCEL,
	}, {
		name:     "cel, bare dollar with no identifier is left as-is",
		in:       "(cel;$ == 'bar')",
		wantStmt: "$ == 'bar'",
		wantEng:  expression.CompilerCEL,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseExpressionRegex(context.TODO(), tt.in)
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
			if assert.NotNil(t, got) {
				assert.Equal(t, tt.wantStmt, got.Statement)
				assert.Equal(t, tt.wantEng, got.Engine)
			}
		})
	}
}
