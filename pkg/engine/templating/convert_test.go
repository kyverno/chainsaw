package templating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convert(t *testing.T) {
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
