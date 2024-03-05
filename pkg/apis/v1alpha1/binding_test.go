package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinding_CheckName(t *testing.T) {
	tests := []struct {
		name         string
		bindingName  string
		bindingValue Any
		wantErr      bool
	}{{
		name:    "empty",
		wantErr: true,
	}, {
		name:        "simple",
		bindingName: "simple",
		wantErr:     false,
	}, {
		name:        "with dollar",
		bindingName: "$simple",
		wantErr:     true,
	}, {
		name:        "with space",
		bindingName: "simple one",
		wantErr:     true,
	}, {
		name:        "with dot",
		bindingName: "simple.one",
		wantErr:     true,
	}, {
		name:        "forbidden",
		bindingName: "config",
		wantErr:     true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Binding{
				Name:  tt.bindingName,
				Value: tt.bindingValue,
			}
			err := b.CheckName()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBinding_CheckEnvName(t *testing.T) {
	tests := []struct {
		name         string
		bindingName  string
		bindingValue Any
		wantErr      bool
	}{{
		name:    "empty",
		wantErr: true,
	}, {
		name:        "simple",
		bindingName: "simple",
		wantErr:     false,
	}, {
		name:        "with dollar",
		bindingName: "$simple",
		wantErr:     true,
	}, {
		name:        "with space",
		bindingName: "simple one",
		wantErr:     true,
	}, {
		name:        "with dot",
		bindingName: "simple.one",
		wantErr:     true,
	}, {
		name:        "forbidden",
		bindingName: "config",
		wantErr:     false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Binding{
				Name:  tt.bindingName,
				Value: tt.bindingValue,
			}
			err := b.CheckEnvName()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
