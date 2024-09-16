package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jpEnv(t *testing.T) {
	tests := []struct {
		name      string
		arguments []any
		want      any
		wantErr   bool
	}{{
		name:      "nil",
		arguments: nil,
		want:      nil,
		wantErr:   true,
	}, {
		name:      "empty",
		arguments: []any{},
		want:      nil,
		wantErr:   true,
	}, {
		name:      "not found",
		arguments: []any{"FOO"},
		want:      "",
		wantErr:   false,
	}, {
		name:      "found",
		arguments: []any{"BAR"},
		want:      "some value",
		wantErr:   false,
	}, {
		name:      "wrong type",
		arguments: []any{12},
		want:      nil,
		wantErr:   true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BAR", "some value")
			got, err := jpEnv(tt.arguments)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
