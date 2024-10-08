package expressions

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
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
		in:       ptr.To(""),
		bindings: apis.NewBindings(),
		want:     ptr.To(""),
		wantErr:  false,
	}, {
		name:     "null",
		in:       ptr.To("(null)"),
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  false,
	}, {
		name:     "error",
		in:       ptr.To("($foo)"),
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  true,
	}, {
		name:     "not string",
		in:       ptr.To("(`42`)"),
		bindings: apis.NewBindings(),
		want:     nil,
		wantErr:  true,
	}, {
		name:     "string",
		in:       ptr.To("('foo')"),
		bindings: apis.NewBindings(),
		want:     ptr.To("foo"),
		wantErr:  false,
	}, {
		name:     "string",
		in:       ptr.To("foo"),
		bindings: apis.NewBindings(),
		want:     ptr.To("foo"),
		wantErr:  false,
	}, {
		name:     "binding",
		in:       ptr.To("($foo)"),
		bindings: apis.NewBindings().Register("$foo", apis.NewBinding("bar")),
		want:     ptr.To("bar"),
		wantErr:  false,
	}, {
		name:     "binding",
		in:       ptr.To("($foo)"),
		bindings: apis.NewBindings().Register("$foo", apis.NewBinding(ptr.To("bar"))),
		want:     ptr.To("bar"),
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
