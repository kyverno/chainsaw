package mutate

import (
	"context"
	"testing"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/stretchr/testify/assert"
)

// ctx, nil, mutation.Parse(ctx, variable.Value.Value), input, bindings, template.WithFunctionCaller(functions.Caller))
func TestMutate(t *testing.T) {
	tests := []struct {
		name     string
		mutation Mutation
		value    any
		bindings binding.Bindings
		want     any
		wantErr  bool
	}{{
		name:     "nil",
		mutation: Parse(context.TODO(), nil),
		value:    42,
		bindings: nil,
		want:     nil,
		wantErr:  false,
	}, {
		name:     "42",
		mutation: Parse(context.TODO(), 64),
		value:    42,
		bindings: nil,
		want:     64,
		wantErr:  false,
	}, {
		name:     "abc",
		mutation: Parse(context.TODO(), "('abc')"),
		value:    "42",
		bindings: nil,
		want:     "abc",
		wantErr:  false,
	}, {
		name: "add",
		mutation: Parse(context.TODO(), map[string]any{
			"c": "(a+b)",
		}),
		value: map[string]any{
			"c": map[string]any{
				"a": 12,
				"b": 24,
			},
		},
		bindings: nil,
		want: map[any]any{
			"c": 36.0,
		},
		wantErr: false,
	}, {
		name: "add",
		mutation: Parse(context.TODO(), map[string]any{
			"c": []any{"(a+b)"},
		}),
		value: map[string]any{
			"c": []any{
				map[string]any{
					"a": 12,
					"b": 24,
				},
			},
		},
		bindings: nil,
		want: map[any]any{
			"c": []any{36.0},
		},
		wantErr: false,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mutate(context.TODO(), nil, tt.mutation, tt.value, tt.bindings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
