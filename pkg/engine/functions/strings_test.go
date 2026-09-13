package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_jpTrimSpace(t *testing.T) {
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
		name:      "wrong type",
		arguments: []any{12},
		want:      nil,
		wantErr:   true,
	}, {
		name:      "no whitespace",
		arguments: []any{"foo"},
		want:      "foo",
		wantErr:   false,
	}, {
		name:      "leading and trailing whitespace",
		arguments: []any{"  foo bar  "},
		want:      "foo bar",
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jpTrimSpace(tt.arguments)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
