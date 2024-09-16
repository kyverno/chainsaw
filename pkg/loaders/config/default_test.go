package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		_fs     func() ([]byte, error)
		wantErr bool
	}{{
		name:    "default",
		wantErr: false,
	}, {
		name: "error",
		_fs: func() ([]byte, error) {
			return nil, errors.New("dummy")
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultConfiguration(tt._fs)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestDefaultConfiguration(t *testing.T) {
	data, err := DefaultConfiguration()
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
