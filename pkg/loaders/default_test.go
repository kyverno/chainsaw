package loaders

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultLoader(t *testing.T) {
	tests := []struct {
		name    string
		_fs     func() (fs.FS, error)
		wantErr bool
	}{{
		name:    "default",
		wantErr: false,
	}, {
		name: "error",
		_fs: func() (fs.FS, error) {
			return nil, errors.New("dummy")
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultLoader(tt._fs)
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

func TestDefaultLoader(t *testing.T) {
	data, err := DefaultLoader()
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
