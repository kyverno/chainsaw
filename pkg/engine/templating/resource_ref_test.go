package templating

import (
	"context"
	"testing"

	"github.com/kyverno/chainsaw/pkg/apis"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestResourceRef(t *testing.T) {
	var noMeta unstructured.Unstructured
	noMeta.SetAPIVersion("v1")
	noMeta.SetKind("Namespace")
	var meta unstructured.Unstructured
	meta.SetAPIVersion("v1")
	meta.SetKind("Namespace")
	meta.SetName("foo")
	meta.SetNamespace("bar")
	meta.SetLabels(
		map[string]string{
			"($foo)": "bar",
		},
	)
	var binds unstructured.Unstructured
	binds.SetAPIVersion("v1")
	binds.SetKind("Namespace")
	binds.SetName("($foo)")
	binds.SetNamespace("($bar)")
	binds.SetLabels(
		map[string]string{
			"($foo)": "($bar)",
		},
	)
	tests := []struct {
		name     string
		obj      *unstructured.Unstructured
		bindings apis.Bindings
		wantErr  bool
		want     *unstructured.Unstructured
	}{{
		name:     "nil",
		obj:      nil,
		bindings: apis.NewBindings(),
		wantErr:  false,
		want:     nil,
	}, {
		name:     "empty",
		obj:      &unstructured.Unstructured{},
		bindings: apis.NewBindings(),
		wantErr:  false,
		want:     &unstructured.Unstructured{},
	}, {
		name:     "meta",
		obj:      noMeta.DeepCopy(),
		bindings: apis.NewBindings(),
		wantErr:  false,
		want:     &noMeta,
	}, {
		name:     "no meta",
		obj:      meta.DeepCopy(),
		bindings: apis.NewBindings(),
		wantErr:  false,
		want:     &meta,
	}, {
		name:     "bindings",
		obj:      binds.DeepCopy(),
		bindings: apis.NewBindings().Register("$foo", apis.NewBinding("foo")).Register("$bar", apis.NewBinding("bar")),
		wantErr:  false,
		want:     &meta,
	}, {
		name:     "error",
		obj:      binds.DeepCopy(),
		bindings: apis.NewBindings(),
		wantErr:  true,
		want:     &binds,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ResourceRef(context.TODO(), apis.XDefaultCompilers, tt.obj, tt.bindings)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, tt.obj)
		})
	}
}
