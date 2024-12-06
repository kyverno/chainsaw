package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinding_CheckName(t *testing.T) {
	tests := []struct {
		name         string
		bindingName  Expression
		bindingValue Projection
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
		name:        "good expression",
		bindingName: "('test')",
		wantErr:     false,
	}, {
		name:        "bad expression",
		bindingName: "('test'",
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
