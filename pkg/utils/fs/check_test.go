package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFolders(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{{
		name:    "nil",
		paths:   nil,
		wantErr: false,
	}, {
		name:    "empty",
		paths:   []string{},
		wantErr: false,
	}, {
		name: "valid",
		paths: []string{
			".",
			"..",
		},
		wantErr: false,
	}, {
		name: "invalid",
		paths: []string{
			".",
			"../foo",
		},
		wantErr: true,
	}, {
		name: "invalid",
		paths: []string{
			"../foo",
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckFolders(tt.paths...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
