package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *rest.Config
		wantErr bool
	}{{
		name:    "nil config",
		cfg:     nil,
		wantErr: true,
	}, {
		name: "valid config",
		cfg: &rest.Config{
			Host:    "http://localhost",
			APIPath: "/api",
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.cfg)
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
