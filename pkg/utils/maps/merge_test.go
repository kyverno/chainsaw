package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name string
		a    map[string]any
		b    map[string]any
		want map[string]any
	}{{
		name: "both null",
		want: map[string]any{},
	}, {
		name: "a null",
		b:    map[string]any{},
		want: map[string]any{},
	}, {
		name: "b null",
		a:    map[string]any{},
		want: map[string]any{},
	}, {
		name: "both empty",
		a:    map[string]any{},
		b:    map[string]any{},
		want: map[string]any{},
	}, {
		name: "",
		a: map[string]any{
			"foo": "bar",
		},
		b: map[string]any{
			"foo": 42,
		},
		want: map[string]any{
			"foo": 42,
		},
	}, {
		name: "",
		a: map[string]any{
			"foo": map[string]any{
				"bar": 42,
			},
		},
		b: map[string]any{
			"foo": nil,
		},
		want: map[string]any{
			"foo": nil,
		},
	}, {
		name: "",
		a: map[string]any{
			"foo": map[string]any{
				"bar": 42,
			},
		},
		b: map[string]any{
			"foo": map[string]any{
				"bar": nil,
			},
		},
		want: map[string]any{
			"foo": map[string]any{
				"bar": nil,
			},
		},
	}, {
		name: "",
		a: map[string]any{
			"foo": 42,
		},
		b: map[string]any{
			"foo": map[string]any{
				"bar": nil,
			},
		},
		want: map[string]any{
			"foo": map[string]any{
				"bar": nil,
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
