package templating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertMap(t *testing.T) {
	tests := []struct {
		name string
		in   any
		out  map[string]any
	}{{
		name: "nil",
	}, {
		name: "int",
		in:   42,
		out:  nil,
	}, {
		name: "ok",
		in: map[any]any{
			"foo": "bar",
		},
		out: map[string]any{
			"foo": "bar",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertMap(tt.in)
			assert.Equal(t, tt.out, got)
		})
	}
}

func Test_convertSlice(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want []any
	}{{
		name: "nil",
	}, {
		name: "int",
		in:   42,
		want: nil,
	}, {
		name: "ok",
		in: []any{
			"foo",
			"bar",
		},
		want: []any{
			"foo",
			"bar",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertSlice(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_convert(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{{
		name: "nil",
	}, {
		name: "int",
		in:   42,
		want: 42,
	}, {
		name: "slice",
		in: []any{
			"foo",
			"bar",
		},
		want: []any{
			"foo",
			"bar",
		},
	}, {
		name: "ok",
		in: map[any]any{
			"foo": "bar",
		},
		want: map[string]any{
			"foo": "bar",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
