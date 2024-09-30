package checks

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name     string
		obj      any
		bindings apis.Bindings
		check    *v1alpha1.Check
		want     field.ErrorList
		wantErr  bool
	}{{
		name:     "nil check",
		obj:      nil,
		bindings: nil,
		check:    nil,
		want:     nil,
		wantErr:  true,
	}, {
		name:     "nil check value",
		obj:      nil,
		bindings: nil,
		check:    &v1alpha1.Check{},
		want:     nil,
		wantErr:  true,
	}, {
		name: "passing",
		obj: map[string]any{
			"foo": "bar",
		},
		bindings: nil,
		check: ptr.To(v1alpha1.NewCheck(
			map[string]any{
				"foo": "bar",
			},
		)),
		want:    nil,
		wantErr: false,
	}, {
		name: "not passing",
		obj: map[string]any{
			"foo": "bar",
		},
		bindings: nil,
		check: ptr.To(v1alpha1.NewCheck(
			map[string]any{
				"foo": "baz",
			},
		)),
		want: []*field.Error{{
			Type:     field.ErrorTypeInvalid,
			Field:    "foo",
			BadValue: "bar",
			Detail:   "Expected value: \"baz\"",
		}},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Check(context.TODO(), apis.DefaultCompilers, tt.obj, tt.bindings, tt.check)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
