package expressions

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		bindings apis.Bindings
		want     string
		wantErr  bool
	}{{
		name:     "empty",
		in:       "",
		bindings: apis.NewBindings(),
		want:     "",
		wantErr:  false,
	}, {
		name:     "error",
		in:       "($foo)",
		bindings: apis.NewBindings(),
		want:     "",
		wantErr:  true,
	}, {
		name:     "not string",
		in:       "(`42`)",
		bindings: apis.NewBindings(),
		want:     "",
		wantErr:  true,
	}, {
		name:     "string",
		in:       "('foo')",
		bindings: apis.NewBindings(),
		want:     "foo",
		wantErr:  false,
	}, {
		name:     "string",
		in:       "foo",
		bindings: apis.NewBindings(),
		want:     "foo",
		wantErr:  false,
	}, {
		name:     "binding",
		in:       "($foo)",
		bindings: apis.NewBindings().Register("$foo", apis.NewBinding("bar")),
		want:     "bar",
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := String(context.TODO(), apis.DefaultCompilers, tt.in, tt.bindings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
