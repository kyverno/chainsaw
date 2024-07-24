package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		in   []string
		want []string
	}{{
		name: "nil",
		env:  nil,
		in:   []string{"echo", "$NAMESPACE"},
		want: []string{"echo", "$NAMESPACE"},
	}, {
		name: "empty",
		env:  map[string]string{},
		in:   []string{"echo", "$NAMESPACE"},
		want: []string{"echo", "$NAMESPACE"},
	}, {
		name: "expand",
		env:  map[string]string{"NAMESPACE": "foo"},
		in:   []string{"echo", "$NAMESPACE"},
		want: []string{"echo", "foo"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Expand(tt.env, tt.in...)
			assert.Equal(t, tt.want, got)
		})
	}
}
