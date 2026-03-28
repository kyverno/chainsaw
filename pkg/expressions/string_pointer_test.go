package expressions

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/stretchr/testify/assert"
)

func TestStringPointer(t *testing.T) {
	tests := []struct {
		name     string
		in       *string
		bindings apis.Bindings
		want     *string
		wantErr  bool
	}{{
		name:     "nil",
		in:       nil,
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  false,
	}, {
		name:     "empty",
		in:       new(""),
		bindings: apis.NewBindings(),
		want:     new(""),
		wantErr:  false,
	}, {
		name:     "null",
		in:       new("(null)"),
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  false,
	}, {
		name:     "error",
		in:       new("($foo)"),
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  true,
	}, {
		name:     "not string",
		in:       new("(`42`)"),
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  true,
	}, {
		name:     "string",
		in:       new("('foo')"),
		bindings: apis.NewBindings(),
		want:     new("foo"),
		wantErr:  false,
	}, {
		name:     "string",
		in:       new("foo"),
		bindings: apis.NewBindings(),
		want:     new("foo"),
		wantErr:  false,
	}, {
		name:     "binding",
		in:       new("($foo)"),
		bindings: apis.NewBindings().Register("$foo", apis.NewBinding("bar")),
		want:     new("bar"),
		wantErr:  false,
	}, {
		name:     "binding",
		in:       new("($foo)"),
		bindings: apis.NewBindings().Register("$foo", apis.NewBinding(new("bar"))),
		want:     new("bar"),
		wantErr:  false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringPointer(context.TODO(), apis.DefaultCompilers, tt.in, tt.bindings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
